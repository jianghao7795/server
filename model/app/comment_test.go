package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommentStruct(t *testing.T) {
	t.Run("验证Comment结构体字段", func(t *testing.T) {
		comment := Comment{
			PostId:   1,
			ParentId: 0,
			Content:  "test comment",
			UserId:   100,
			ToUserId: 200,
		}
		assert.Equal(t, uint(1), comment.PostId)
		assert.Equal(t, uint(0), comment.ParentId)
		assert.Equal(t, "test comment", comment.Content)
		assert.Equal(t, uint(100), comment.UserId)
		assert.Equal(t, uint(200), comment.ToUserId)
	})

	t.Run("验证表名", func(t *testing.T) {
		assert.Equal(t, "comments", Comment{}.TableName())
	})
}

func TestCommentChildrenField(t *testing.T) {
	comment := Comment{
		Content: "parent comment",
		Children: []Comment{
			{Content: "child 1"},
			{Content: "child 2"},
		},
	}
	assert.Len(t, comment.Children, 2)
	assert.Equal(t, "child 1", comment.Children[0].Content)
	assert.Equal(t, "child 2", comment.Children[1].Content)
}
