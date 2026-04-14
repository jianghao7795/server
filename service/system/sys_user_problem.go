package system

import (
	global "server/model"
	"server/model/system"
	"server/model/system/response"
	"strconv"
	"strings"
)

func (*Problem) GetUserProblemList(info *system.SysUserProblem) (list any, err error) {
	db := global.DB.Model(&system.SysUserProblem{})
	var problemList []response.SysUserProblem
	err = db.Where("sys_user_id = ?", info.SysUserId).Find(&problemList).Error
	return problemList, err
}

func (*Problem) GetUserProblemSettingList(info *system.SysUserProblem) (list any, err error) {
	db := global.DB.Model(&system.SysUserProblem{})
	var problemList []response.SysUserProblemSetting
	err = db.Where("sys_user_id = ?", info.SysUserId).Find(&problemList).Error
	return problemList, err
}

func (*Problem) SetUserProblemSetting(problem []system.SysUserProblem) (string, error) {
	if len(problem) == 0 {
		return "没有新建和更新", nil
	}

	var messageString strings.Builder
	for index, item := range problem {
		if item.ID == 0 {
			db := global.DB.Model(&system.SysUserProblem{})
			err := db.Create(&item).Error
			if err != nil {
				return "", err
			}
			messageString.WriteString("[" + strconv.Itoa(index) + "]" + "新建成功 ")
		} else {
			db := global.DB.Model(&system.SysUserProblem{})
			var dataProblemFirst system.SysUserProblem
			err := db.Where("id = ?", item.ID).First(&dataProblemFirst).Error
			if err != nil {
				return "", err
			}
			dataProblemFirst.Answer = item.Answer
			dataProblemFirst.Problem = item.Problem
			dataProblemFirst.SysUserId = item.SysUserId
			err = db.Save(&dataProblemFirst).Error
			if err != nil {
				return "", err
			}
			messageString.WriteString("[" + strconv.Itoa(index) + "]" + "更新成功 ")
		}
	}

	return messageString.String(), nil
}

func (*Problem) HasSetting(uid int) (bool, error) {
	if uid == 0 {
		return false, nil
	}
	db := global.DB.Model(&system.SysUserProblem{})
	var total int64
	err := db.Where("sys_user_id = ?", uid).Order("id desc").Count(&total).Error
	if err != nil {
		return false, err
	}
	return total != 0, nil
}

func (*Problem) VerifyAnswer(info *system.SysUserProblem) (bool, error) {
	db := global.DB.Model(&system.SysUserProblem{})
	var findProblem system.SysUserProblem
	err := db.Where("id = ?", info.ID).First(&findProblem).Error
	if err != nil {
		return false, err
	}
	isAnswer := findProblem.Answer == info.Answer
	return isAnswer, nil
}
