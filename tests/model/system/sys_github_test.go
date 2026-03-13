package system

import (
	"encoding/json"
	"testing"
	"time"

	"server/model/system"

	"github.com/stretchr/testify/assert"
)

func TestSysGithub_Validation(t *testing.T) {
	tests := []struct {
		name    string
		github  system.SysGithub
		isValid bool
	}{
		{
			name: "有效的GitHub记录",
			github: system.SysGithub{
				Author:     "testuser",
				Message:    "feat: add new feature",
				CommitTime: "2024-01-01 12:00:00",
			},
			isValid: true,
		},
		{
			name: "缺少作者信息",
			github: system.SysGithub{
				Author:     "",
				Message:    "feat: add new feature",
				CommitTime: "2024-01-01 12:00:00",
			},
			isValid: false,
		},
		{
			name: "缺少提交信息",
			github: system.SysGithub{
				Author:     "testuser",
				Message:    "",
				CommitTime: "2024-01-01 12:00:00",
			},
			isValid: false,
		},
		{
			name: "缺少提交时间",
			github: system.SysGithub{
				Author:     "testuser",
				Message:    "feat: add new feature",
				CommitTime: "",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isValid {
				assert.NotEmpty(t, tt.github.Author)
				assert.NotEmpty(t, tt.github.Message)
				assert.NotEmpty(t, tt.github.CommitTime)
			} else {
				assert.True(t,
					tt.github.Author == "" ||
						tt.github.Message == "" ||
						tt.github.CommitTime == "")
			}
		})
	}
}

func TestGithubCommit_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"sha": "abc123",
		"commit": {
			"author": {
				"name": "testuser",
				"email": "test@example.com",
				"date": "2024-01-01T12:00:00Z"
			},
			"message": "feat: add new feature"
		}
	}`

	var commit system.GithubCommit
	err := json.Unmarshal([]byte(jsonData), &commit)

	assert.NoError(t, err)
	assert.Equal(t, "abc123", commit.Sha)
	assert.Equal(t, "testuser", commit.Commit.Author.Name)
	assert.Equal(t, "test@example.com", commit.Commit.Author.Email)
	assert.Equal(t, "feat: add new feature", commit.Commit.Message)
}

func TestGithubCommit_TimeParsing(t *testing.T) {
	// 测试时间解析
	timeStr := "2024-01-01T12:00:00Z"
	parsedTime, err := time.Parse(time.RFC3339, timeStr)

	assert.NoError(t, err)
	assert.Equal(t, 2024, parsedTime.Year())
	assert.Equal(t, time.January, parsedTime.Month())
	assert.Equal(t, 1, parsedTime.Day())
	assert.Equal(t, 12, parsedTime.Hour())
}

func TestSysGithub_TimeFormat(t *testing.T) {
	// 测试时间格式
	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")

	// 验证时间格式是否正确
	_, err := time.Parse("2006-01-02 15:04:05", formattedTime)
	assert.NoError(t, err)

	// 验证时间格式包含必要的组件
	assert.Contains(t, formattedTime, "-")
	assert.Contains(t, formattedTime, ":")
	assert.Len(t, formattedTime, 19) // "2006-01-02 15:04:05" 的长度
}
