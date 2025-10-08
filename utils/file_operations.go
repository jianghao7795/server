package utils

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

//@function: FileMove
//@description: 文件移动供外部调用
//@param: src string, dst string(src: 源位置,绝对路径or相对路径, dst: 目标位置,绝对路径or相对路径,必须为文件夹)
//@return: err error

func FileMove(src string, dst string) (err error) {
	if dst == "" {
		return nil
	}
	src, err = filepath.Abs(src)
	if err != nil {
		return err
	}
	dst, err = filepath.Abs(dst)
	if err != nil {
		return err
	}
	revoke := false
	dir := filepath.Dir(dst)
Redirect:
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			return err
		}
		if !revoke {
			revoke = true
			goto Redirect
		}
	}
	return os.Rename(src, dst)
}

// DeleteFile removes a file or directory
// @param filePath path to the file or directory to delete
// @return error if deletion fails
func DeleteFile(filePath string) error {
	return os.RemoveAll(filePath)
}

//@function: TrimSpace
//@description: 去除结构体空格
//@param: target interface (target: 目标结构体,传入必须是指针类型)
//@return: null

func TrimSpace(target any) {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Ptr {
		return
	}
	t = t.Elem()
	v := reflect.ValueOf(target).Elem()
	for i := range t.NumField() {
		switch v.Field(i).Kind() {
		case reflect.String:
			v.Field(i).SetString(strings.TrimSpace(v.Field(i).String()))
		}
	}
}

// @function: FileExist
// @description: 判断文件是否存在
// @param: path string (path: 文件路径)
// @return: bool (true: 文件存在, false: 文件不存在)
// FileExist 判断文件是否存在
func FileExist(path string) bool {
	fi, err := os.Lstat(path)
	if err == nil {
		return !fi.IsDir()
	}
	return !os.IsNotExist(err)
}

// 读取文件 并按行读取 并执行回调函数
// @function: ReadFile
// @description: 读取文件 并按行读取 并执行回调函数
// @param: path string (path: 文件路径)
// @param: fn func(string) error (fn: 回调函数, 传入每行内容)
// @return: error (nil: 读取成功, 其他: 读取失败)
func ReadFile(path string, fn func(string) error) error {
	const BufferSize = 100 // 缓冲区大小
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	b := make([]byte, BufferSize)
	for {
		n, err := f.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		err = fn(string(b[:n]))
		if err != nil {
			return err
		}

		if n < BufferSize {
			break
		}
	}
	return nil
}
