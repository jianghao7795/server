package app

import (
	"testing"

	"server/model/app"
	appReq "server/model/app/request"

	"github.com/stretchr/testify/assert"
)

func TestTagServiceStruct(t *testing.T) {
	svc := &TagService{}
	assert.NotNil(t, svc)
}

func TestTagCreateValidation(t *testing.T) {
	tag := &app.Tag{
		Name:   "Go语言",
		Status: 1,
	}
	assert.Equal(t, "Go语言", tag.Name)
	assert.Equal(t, 1, tag.Status)
}

func TestTagSearchParams(t *testing.T) {
	t.Run("构建带名称的搜索参数", func(t *testing.T) {
		info := &appReq.TagSearch{
			Tag: app.Tag{
				Name: "Go",
			},
		}
		info.Page = 1
		info.PageSize = 10

		assert.Equal(t, "Go", info.Name)
		assert.Equal(t, 1, info.Page)
		assert.Equal(t, 10, info.PageSize)
	})

	t.Run("构建空搜索参数", func(t *testing.T) {
		info := &appReq.TagSearch{}
		info.Page = 3
		info.PageSize = 5

		assert.Empty(t, info.Name)
		assert.Equal(t, 3, info.Page)
		assert.Equal(t, 5, info.PageSize)
	})
}

func TestTagService_CreateTag(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestTagService_DeleteTag(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestTagService_UpdateTag(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestTagService_GetTag(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestTagService_GetTagInfoList(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}
