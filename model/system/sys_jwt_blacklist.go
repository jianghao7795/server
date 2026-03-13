package system

import (
	global "server/model"
)

type JwtBlacklist struct {
	global.MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
