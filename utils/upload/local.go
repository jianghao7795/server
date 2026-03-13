package upload

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path"
	"server/utils"
	"strings"
	"time"

	global "server/model"

	"go.uber.org/zap"
)

type Local struct{}

// @object: *Local
// @function: UploadFile
// @description: 本地上传文件
// @param: file *multipart.FileHeader
// @return: string, string, error
func (*Local) UploadFile(file *multipart.FileHeader) (string, string, error) {
	// 读取文件后缀
	ext := path.Ext(file.Filename)
	// 读取文件名并加密
	name := strings.TrimSuffix(file.Filename, ext)
	name = utils.MD5V([]byte(name))

	nowTime := time.Now()
	// 路径
	filepath := nowTime.Format("2006/01/02")
	// 拼接新文件名
	filename := name + "_" + nowTime.Format("150405") + ext
	// 尝试创建此路径
	mkdirErr := os.MkdirAll(global.CONFIG.Local.Path+"/"+filepath, os.ModePerm)
	if mkdirErr != nil {
		global.LOG.Error("function os.MkdirAll() Filed", zap.Any("err", mkdirErr.Error()))
		return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}
	// 拼接路径和文件名
	p := global.CONFIG.Local.Path + "/" + filepath + "/" + filename

	f, openError := file.Open() // 读取文件
	if openError != nil {
		global.LOG.Error("function file.Open() Filed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // 文件 defer 关闭

	out, createErr := os.Create(p)
	if createErr != nil {
		global.LOG.Error("function os.Create() Filed", zap.Any("err", createErr.Error()))
		return "", "", errors.New("function os.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close() // 创建文件 defer 关闭

	_, copyErr := io.Copy(out, f) // 传输（拷贝）文件
	if copyErr != nil {
		global.LOG.Error("function io.Copy() Filed", zap.Any("err", copyErr.Error()))
		return "", "", errors.New("function io.Copy() Filed, err:" + copyErr.Error())
	}
	return p, filename, nil
}

//@object: *Local
//@function: DeleteFile
//@description: 删除文件
//@param: key string
//@return: error

func (*Local) DeleteFile(key string) error {
	if strings.Contains(key, global.CONFIG.Local.Path) {
		_, err := os.Stat(key)
		if err != nil {
			if os.IsNotExist(err) {
				global.LOG.Error("本地文件删除失败, err: 文件不存在 ", zap.Error(err))
				return nil
			} else {
				return err
			}
		}

		if err := os.Remove(key); err != nil {
			return errors.New("本地文件删除失败, err:" + err.Error())
		}
	} else {
		return errors.New("本地文件删除失败, err: 请配置环境是否为本地文件")
	}

	return nil
}
