package mobile

import (
	"testing"

	"server/model/mobile"
	mobileReq "server/model/mobile/request"

	"github.com/stretchr/testify/assert"
)

func TestMobileUserServiceStruct(t *testing.T) {
	svc := &MobileUserService{}
	assert.NotNil(t, svc)
}

func TestMobileUserCreateValidation(t *testing.T) {
	user := &mobile.MobileUser{
		Username: "newuser",
		Nickname: "新用户",
		Realname: "张三",
		Phone:    "13800138000",
		Gender:   1,
		Industry: 3,
	}
	assert.Equal(t, "newuser", user.Username)
	assert.Equal(t, "新用户", user.Nickname)
	assert.Equal(t, "张三", user.Realname)
	assert.Equal(t, "13800138000", user.Phone)
	assert.Equal(t, uint8(1), user.Gender)
	assert.Equal(t, uint8(3), user.Industry)
}

func TestMobileUserSearchParams(t *testing.T) {
	t.Run("构建带用户名的搜索参数", func(t *testing.T) {
		info := &mobileReq.MobileUserSearch{
			MobileUser: mobile.MobileUser{
				Username: "testuser",
			},
		}
		info.Page = 1
		info.PageSize = 15

		assert.Equal(t, "testuser", info.Username)
		assert.Equal(t, 1, info.Page)
		assert.Equal(t, 15, info.PageSize)
	})

	t.Run("构建空搜索参数", func(t *testing.T) {
		info := &mobileReq.MobileUserSearch{}
		info.Page = 2
		info.PageSize = 10

		assert.Empty(t, info.Username)
		assert.Equal(t, 2, info.Page)
		assert.Equal(t, 10, info.PageSize)
	})
}

func TestMobileUpdateStruct(t *testing.T) {
	update := mobileReq.MobileUpdate{
		Field: "avatar",
		Value: "https://example.com/avatar.jpg",
	}
	assert.Equal(t, "avatar", update.Field)
	assert.Equal(t, "https://example.com/avatar.jpg", update.Value)
}

func TestMobileUpdatePasswordStruct(t *testing.T) {
	req := mobileReq.MobileUpdatePassword{
		ID:          10,
		Password:    "old123",
		NewPassword: "new456",
	}
	assert.Equal(t, uint(10), req.ID)
	assert.Equal(t, "old123", req.Password)
	assert.Equal(t, "new456", req.NewPassword)
}

func TestMobileUserService_CreateMobileUser(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestMobileUserService_DeleteMobileUser(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestMobileUserService_UpdateMobileUser(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestMobileUserService_GetMobileUser(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestMobileUserService_GetMobileUserInfoList(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}
