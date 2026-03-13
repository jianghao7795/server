package system

import (
	global "server/model"
	"server/model/common/request"
	"server/model/system"
	"sort"
	"time"
)

func (g *GithubService) CreateGithub(github *[]system.SysGithub) (total int, err error) {
	db := global.DB.Model(&system.SysGithub{})
	var data []system.SysGithub
	for _, item := range *github {
		if item.CommitTime != "" {
			db = db.Or("commit_time = ?", item.CommitTime)
		}
	}
	db.Order("id desc").Find(&data)
	dataInsert := []system.SysGithub{}
	isExist := true
	for _, item := range *github {
		for _, itemGithub := range data {
			if itemGithub.CommitTime == item.CommitTime {
				isExist = false
				break
			} else {
				isExist = true
			}
		}
		if isExist {
			dataInsert = append(dataInsert, item)
			isExist = true
		}
	}
	insertNumber := 0
	if len(dataInsert) != 0 {
		sort.Slice(dataInsert, func(i, j int) bool {
			return timeStrTime(dataInsert[i].CommitTime) < timeStrTime(dataInsert[j].CommitTime)
		})
		for _, item := range dataInsert {
			dataItem := item
			if dataItem.CommitTime != "" {
				insertNumber++
				err = db.Create(&dataItem).Error
			}
		}
	}

	return insertNumber, err
}

func (g *GithubService) GetGithubList(info request.PageInfo) (list []system.SysGithub, err error) {
	db := global.DB.Model(&system.SysGithub{})
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&list).Error
	return
}

func timeStrTime(valueStr string) int64 {
	loc := time.Local
	fmtStr := "2006-01-02 15:04:05"
	t, _ := time.ParseInLocation(fmtStr, valueStr, loc)
	return t.Unix()
}
