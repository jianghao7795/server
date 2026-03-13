package example

import (
	"mime/multipart"
	global "server/model"
)

// file struct, 文件结构体
type ExaFile struct {
	global.MODEL
	FileName     string         `json:"file_name" query:"file_name" form:"file_name" gorm:"comment:文件名"`
	FileMd5      string         `json:"file_md5" query:"file_md5" form:"file_md5" gorm:"comment:文件的md5"`
	FilePath     string         `json:"file_path" query:"file_path" form:"file_path" gorm:"comment:文件路径"`
	ExaFileChunk []ExaFileChunk `json:"exa_file_chunk" query:"exa_file_chunk" form:"exa_file_chunk"` // 切片结构体
	ChunkTotal   int            `json:"chunk_total" query:"chunk_total" form:"chunk_total" gorm:"comment:切片总数"`
	IsFinish     bool           `json:"is_finish" query:"is_finish" form:"is_finish" gorm:"comment:是否完成"`
}

func (ExaFile) TableName() string {
	return "exa_files"
}

// file chunk struct, 切片结构体
type ExaFileChunk struct {
	global.MODEL
	ExaFileID       uint
	FileChunkNumber int
	FileChunkPath   string
}

func (ExaFileChunk) TableName() string {
	return "exa_file_chunks"
}

type ExaFileData struct {
	FileName string `json:"fileName" form:"fileName"`
	FileMd5  string `json:"fileMd5" form:"fileMd5"`
	// FilePath string `json:"file_path" form:"file_path"`
	ChunkMd5    string                `json:"chunkMd5" form:"chunkMd5"`
	ChunkNumber int                   `json:"chunkNumber" form:"chunkNumber"`
	ChunkTotal  int                   `json:"chunkTotal" form:"chunkTotal"`
	FileHeader  *multipart.FileHeader `json:"file" form:"file"`
}
