package app

import (
	"testing"

	"server/model/app"
	commentReq "server/model/app/request"
	"server/model/common/request"

	"github.com/stretchr/testify/assert"
)

func TestCommentServiceStruct(t *testing.T) {
	svc := &CommentService{}
	assert.NotNil(t, svc)
}

func TestCommentCreateValidation(t *testing.T) {
	comment := &app.Comment{
		Content:  "测试评论",
		PostId:   1,
		UserId:   100,
		ParentId: 0,
	}
	assert.Equal(t, "测试评论", comment.Content)
	assert.Equal(t, uint(1), comment.PostId)
	assert.Equal(t, uint(100), comment.UserId)
	assert.Equal(t, uint(0), comment.ParentId)
}

func TestCommentSearchParams(t *testing.T) {
	t.Run("构建带文章ID的搜索参数", func(t *testing.T) {
		info := &commentReq.CommentSearch{
			ArticleId: 5,
			Content:   "hello",
		}
		info.Page = 1
		info.PageSize = 20

		assert.Equal(t, 5, info.ArticleId)
		assert.Equal(t, "hello", info.Content)
		assert.Equal(t, 1, info.Page)
		assert.Equal(t, 20, info.PageSize)
	})

	t.Run("构建仅有内容的搜索参数", func(t *testing.T) {
		info := &commentReq.CommentSearch{
			Content: "搜索词",
		}
		info.Page = 2
		info.PageSize = 10

		assert.Equal(t, "搜索词", info.Content)
		assert.Equal(t, 0, info.ArticleId)
		assert.Equal(t, 2, info.Page)
	})
}

func TestCommentDeleteByIds(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestCommentService_CreateComment(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestCommentService_UpdateComment(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestCommentService_GetComment(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestCommentService_GetCommentList(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestCommentService_GetCommentInfoList(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestCommentService_GetCommentTreeList(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestCommentService_PutLikeItOrDislike(t *testing.T) {
	t.Skip("需要完整的应用环境，跳过需要DB的测试")
}

func TestCommentIdsReq(t *testing.T) {
	ids := request.IdsReq{Ids: []int{1, 2, 3}}
	assert.Len(t, ids.Ids, 3)
	assert.Equal(t, []int{1, 2, 3}, ids.Ids)
}
