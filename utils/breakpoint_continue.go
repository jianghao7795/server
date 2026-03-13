package utils

import (
	"errors"
	ioutil "io"
	"math"
	"os"
	global "server/model"
	"strconv"
)

// 前端传来文件片与当前片为什么文件的第几片
// 后端拿到以后比较次分片是否上传 或者是否为不完全片
// 前端发送每片多大
// 前端告知是否为最后一片且是否完成

const (
	breakpointDir = "uploads/breakpointDir/"
	finishDir     = "uploads/fileDir/"
)

//@author:wuhao
//@function: BreakPointContinue
//@description: 断点续传
//@param: content []byte, fileName string, contentNumber int, contentTotal int, fileMd5 string
//@return: error, string

func BreakPointContinue(content []byte, fileName string, contentNumber int, contentTotal int, fileMd5 string) (string, error) {
	path := breakpointDir + fileMd5 + "/"
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return path, err
	}
	paths, err := makeFileContent(content, fileName, path, contentNumber)
	return paths, err
}

//@author: wuhao
//@function: CheckMd5
//@description: 检查Md5
//@param: content []byte, chunkMd5 string
//@return: CanUpload bool

func CheckMd5(content []byte, chunkMd5 string) (CanUpload bool) {
	fileMd5 := MD5V(content)
	// log.Println("content: ", fileMd5, chunkMd5)
	if fileMd5 == chunkMd5 {
		return true // 可以继续上传
	} else {
		return false // 切片不完整，废弃
	}
}

//@author: wuhao
//@function: makeFileContent
//@description: 创建切片内容
//@param: content []byte, fileName string, FileDir string, contentNumber int
//@return: error, string

func makeFileContent(content []byte, fileName string, FileDir string, contentNumber int) (string, error) {
	path := FileDir + fileName + "_" + strconv.Itoa(contentNumber)
	f, err := os.Create(path)
	if err != nil {
		return path, err
	} else {
		_, err = f.Write(content)
		if err != nil {
			return path, err
		}
	}
	defer f.Close()
	return path, nil
}

//@author: wuhao
//@function: makeFileContent
//@description: 创建切片文件
//@param: fileName string, FileMd5 string
//@return: error, string

func MakeFile(fileName string, FileMd5 string) (string, error) {
	rd, err := os.ReadDir(breakpointDir + FileMd5)
	if err != nil {
		return finishDir + fileName, err
	}
	// if os.IsExist()
	err = os.MkdirAll(finishDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	fd, err := os.OpenFile(finishDir+fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return finishDir + fileName, err
	}
	defer func(fd *os.File) {
		err := fd.Close()
		if err != nil {
			global.LOG.Error("关闭文件失败" + err.Error())
			err = os.Remove(finishDir + fileName)
			if err != nil {
				global.LOG.Error("删除文件失败" + err.Error())
			}
		}
	}(fd)
	for k := range rd {
		content, _ := os.ReadFile(breakpointDir + FileMd5 + "/" + fileName + "_" + strconv.Itoa(k))
		_, err = fd.Write(content)
		if err != nil {
			// log.Println("log", err.Error())
			err = os.Remove(finishDir + fileName)
			if err != nil {
				return "", err
			}
			return finishDir + fileName, err
		}
	}
	return finishDir + fileName, nil
}

//@author: wuhao
//@function: RemoveChunk
//@description: 移除切片
//@param: FileMd5 string
//@return: error

func RemoveChunk(FileMd5 string) error {
	if FileMd5 == "" {
		return errors.New("路径错误")
	}
	err := os.RemoveAll(breakpointDir + FileMd5) // RemoveAll删除path指定的文件，或目录及它包含的任何下级对象。它会尝试删除所有东西，除非遇到错误并返回。如果path指定的对象不存在，RemoveAll会返回nil而不返回错误。
	return err
}

// 文件分为多个片，每个片的大小为1M，每个片都有一个唯一的编号，编号从1开始，编号为0表示整个文件。
// 客户端上传文件时，会先将文件切分为多个片，并将每个片的编号、大小、MD5值等信息记录在服务器端。
// 当客户端断点续传时，会将已上传的片的编号、大小、MD5值等信息记录在客户端，并将这些信息发送给服务器端。
// 服务器端收到客户端的断点续传信息后，会根据这些信息判断客户端是否已经上传了完整的文件。
// 如果客户端已经上传了完整的文件，则服务器端会将该文件从断点续传目录移动到文件目录，并返回文件完整路径。
// 如果客户端未上传完整的文件，则服务器端会将该文件从断点续传目录删除，并返回错误信息。
func FileSeparateMerge(file_path string, chunk_size ...int) error {
	const chunkSize = 1 << (10 * 2)
	chunk_num := chunkSize
	if len(chunk_size) > 0 {
		chunk_num = chunk_size[0]
	}

	fileInfo, err := os.Stat(file_path)
	if err != nil {
		panic(err)
	}

	num := math.Ceil(float64(fileInfo.Size()) / float64(chunk_num))

	fi, err := os.OpenFile("cbd.mp4", os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	b := make([]byte, chunkSize)
	var i int64 = 1
	for ; i <= int64(num); i++ {
		fi.Seek((i-1)*int64(chunk_num), 0)
		if len(b) > int(fileInfo.Size()-(i-1)*int64(chunk_num)) {
			b = make([]byte, fileInfo.Size()-(i-1)*int64(chunk_num))
		}
		fi.Read(b)

		f, err := os.OpenFile("./"+strconv.Itoa(int(i))+".db", os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		f.Write(b)
		f.Close()
	}
	fi.Close()

	fii, err := os.OpenFile("all.mp4", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	for i := 1; i <= int(num); i++ {
		f, err := os.OpenFile("./"+strconv.Itoa(int(i))+".db", os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		fii.Write(b)
		f.Close()
	}
	return nil
}
