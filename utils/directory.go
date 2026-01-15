package utils

import (
	"errors"
	"os"
	global "server-fiber/model"

	"go.uber.org/zap"
)

//@author: wuhao
//@function: PathExists
//@description: 目录是否存在
//@param: path string
//@return: bool, error

func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return false, errors.New("存在同名文件")
	}
	if os.IsNotExist(err) {
		return false, errors.New("目录不存在")
	}
	return false, err
}

// @author: wuhao
// @function: PathExistsFile
// @description: 文件是否存在
// @param: path string
// @return: bool, error
func PathExistsFile(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, errors.New("文件不存在")
	}

	return true, nil
}

//@author: wuhao
//@function: CreateDir
//@description: 批量创建文件夹
//@param: dirs ...string
//@return: err error

func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			global.LOG.Debug("create directory" + v)
			if err := os.MkdirAll(v, os.ModePerm); err != nil {
				global.LOG.Error("create directory"+v, zap.Any(" error:", err))
				return err
			}
		}
	}
	return err
}
