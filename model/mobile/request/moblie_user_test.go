package request

import (
	"testing"

	"server/model/mobile"

	"github.com/stretchr/testify/assert"
)

func TestMobileUserSearchStruct(t *testing.T) {
	search := MobileUserSearch{
		MobileUser: mobile.MobileUser{
			Username: "testuser",
			Phone:    "13800138000",
		},
	}
	search.Page = 1
	search.PageSize = 10

	assert.Equal(t, "testuser", search.Username)
	assert.Equal(t, "13800138000", search.Phone)
	assert.Equal(t, 1, search.Page)
	assert.Equal(t, 10, search.PageSize)
}

func TestMobileUpdateStruct(t *testing.T) {
	update := MobileUpdate{
		Field: "nickname",
		Value: "新昵称",
	}
	assert.Equal(t, "nickname", update.Field)
	assert.Equal(t, "新昵称", update.Value)
}

func TestMobileUpdatePasswordStruct(t *testing.T) {
	req := MobileUpdatePassword{
		ID:          1,
		Password:    "oldpass",
		NewPassword: "newpass",
	}
	assert.Equal(t, uint(1), req.ID)
	assert.Equal(t, "oldpass", req.Password)
	assert.Equal(t, "newpass", req.NewPassword)
}
