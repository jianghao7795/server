package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"server-fiber/model"
	"server-fiber/model/system"
	"sort"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ClearTable clears database table data based on time interval
// @param db database connection
// @param tableName name of the table to clear
// @param compareField field to compare with time
// @param interval time interval string (e.g., "24h", "7d")
// @return error if operation fails
func ClearTable(db *gorm.DB, tableName string, compareField string, interval string) error {
	if db == nil {
		return errors.New("database connection cannot be empty")
	}

	duration, err := time.ParseDuration(interval)
	if err != nil {
		return fmt.Errorf("invalid duration format '%s': %w", interval, err)
	}

	if duration < 0 {
		return errors.New("duration cannot be negative")
	}

	cutoffTime := time.Now().Add(-duration)
	query := fmt.Sprintf("UPDATE %s SET deleted_at = ? WHERE %s < ?", tableName, compareField)

	return db.Exec(query, cutoffTime, cutoffTime).Error
}

// UpdateTable fetches GitHub commits and updates database
// @param db database connection
// @param tableName name of the table (unused in current implementation)
// @param compareField field to compare with time (unused in current implementation)
// @param interval time interval string (unused in current implementation)
// @return error if operation fails
func UpdateTable(db *gorm.DB, tableName string, compareField string, interval string) error {
	if db == nil {
		return errors.New("database connection cannot be empty")
	}

	// Validate parameters (even though not used in current implementation)
	_, err := time.ParseDuration(interval)
	if err != nil {
		return fmt.Errorf("invalid duration format '%s': %w", interval, err)
	}

	// Fetch GitHub commits
	commits, err := fetchGitHubCommits()
	if err != nil {
		model.LOG.Error("Failed to fetch GitHub commits", zap.Error(err))
		return fmt.Errorf("failed to fetch GitHub commits: %w", err)
	}

	// Convert to database format
	githubData := convertCommitsToGithubData(commits)

	// Get existing records
	existingCommits, err := getExistingCommits(db)
	if err != nil {
		model.LOG.Error("Failed to get existing commits", zap.Error(err))
		return fmt.Errorf("failed to get existing commits: %w", err)
	}

	// Find new commits to insert
	newCommits := findNewCommits(githubData, existingCommits)

	// Insert new commits
	insertCount, err := insertNewCommits(db, newCommits)
	if err != nil {
		model.LOG.Error("Failed to insert new commits", zap.Error(err))
		return fmt.Errorf("failed to insert new commits: %w", err)
	}

	model.LOG.Info(fmt.Sprintf("Successfully inserted %d new GitHub commits", insertCount))
	return nil
}

// fetchGitHubCommits fetches commits from GitHub API
func fetchGitHubCommits() ([]system.GithubCommit, error) {
	const (
		page    = "1"
		perPage = "20"
		url     = "https://api.github.com/repos/JiangHaoCode/server-fiber/commits"
	)

	reqURL := fmt.Sprintf("%s?page=%s&per_page=%s", url, page, perPage)
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var commits []system.GithubCommit
	if err := json.Unmarshal(body, &commits); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return commits, nil
}

// convertCommitsToGithubData converts GitHub commits to database format
func convertCommitsToGithubData(commits []system.GithubCommit) []system.SysGithub {
	// Load Shanghai timezone
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.UTC // fallback to UTC
	}

	data := make([]system.SysGithub, 0, len(commits))
	for _, commit := range commits {
		githubData := system.SysGithub{
			Author:     commit.Commit.Author.Name,
			CommitTime: commit.Commit.Author.Date.In(loc).Format("2006-01-02 15:04:05"),
			Message:    commit.Commit.Message,
		}
		data = append(data, githubData)
	}
	return data
}

// getExistingCommits retrieves existing commits from database
func getExistingCommits(db *gorm.DB) ([]system.SysGithub, error) {
	var existingCommits []system.SysGithub
	err := db.Model(&system.SysGithub{}).
		Limit(20).
		Order("id desc").
		Find(&existingCommits).Error
	return existingCommits, err
}

// findNewCommits finds commits that don't exist in database
func findNewCommits(newCommits, existingCommits []system.SysGithub) []system.SysGithub {
	// Create a map for faster lookup
	existingMap := make(map[string]bool)
	for _, existing := range existingCommits {
		existingMap[existing.CommitTime] = true
	}

	var newCommitsToInsert []system.SysGithub
	for _, commit := range newCommits {
		if commit.CommitTime != "" && !existingMap[commit.CommitTime] {
			newCommitsToInsert = append(newCommitsToInsert, commit)
		}
	}

	// Sort by commit time
	sort.Slice(newCommitsToInsert, func(i, j int) bool {
		return timeStrToUnix(newCommitsToInsert[i].CommitTime) < timeStrToUnix(newCommitsToInsert[j].CommitTime)
	})

	return newCommitsToInsert
}

// insertNewCommits inserts new commits into database
func insertNewCommits(db *gorm.DB, commits []system.SysGithub) (int, error) {
	if len(commits) == 0 {
		return 0, nil
	}

	insertCount := 0
	for _, commit := range commits {
		if commit.CommitTime != "" {
			if err := db.Create(&commit).Error; err != nil {
				return insertCount, fmt.Errorf("failed to insert commit %s: %w", commit.CommitTime, err)
			}
			insertCount++
		}
	}

	return insertCount, nil
}

// timeStrToUnix converts time string to Unix timestamp
// @param valueStr time string in format "2006-01-02 15:04:05"
// @return Unix timestamp
func timeStrToUnix(valueStr string) int64 {
	const timeFormat = "2006-01-02 15:04:05"
	loc := time.Local
	t, err := time.ParseInLocation(timeFormat, valueStr, loc)
	if err != nil {
		// Return 0 for invalid time strings
		return 0
	}
	return t.Unix()
}
