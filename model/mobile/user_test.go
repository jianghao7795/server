package mobile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMobileUserStruct(t *testing.T) {
	user := MobileUser{
		Username: "testuser",
		Nickname: "测试昵称",
		Realname: "测试用户",
		Avatar:   "https://avatar.example.com/1.jpg",
		Sign:     "Hello World",
		Phone:    "13800138000",
		Gender:   1,
		Industry: 2,
	}
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "测试昵称", user.Nickname)
	assert.Equal(t, "测试用户", user.Realname)
	assert.Equal(t, "https://avatar.example.com/1.jpg", user.Avatar)
	assert.Equal(t, "Hello World", user.Sign)
	assert.Equal(t, "13800138000", user.Phone)
	assert.Equal(t, uint8(1), user.Gender)
	assert.Equal(t, uint8(2), user.Industry)
}

func TestMobileUserTableName(t *testing.T) {
	assert.Equal(t, "mobile_users", MobileUser{}.TableName())
}
