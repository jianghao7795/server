package model

import (
	"time"

	"gorm.io/gorm"
)

type MODEL struct {
	ID        uint           `gorm:"primarykey"`     // 主键ID
	CreatedAt time.Time      `json:"createdAt"`      // 创建时间
	UpdatedAt time.Time      `json:"updatedAt"`      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

type SearchValue struct {
	Name       string
	Value      string
	SearchType string
}
