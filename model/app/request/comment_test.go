package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommentSearch(t *testing.T) {
	search := CommentSearch{
		ArticleId: 1,
		Content:   "hello",
	}
	assert.Equal(t, 1, search.ArticleId)
	assert.Equal(t, "hello", search.Content)
}

func TestCommentSearchWithPagination(t *testing.T) {
	search := CommentSearch{
		ArticleId: 2,
		Content:   "test",
	}
	search.Page = 1
	search.PageSize = 20

	assert.Equal(t, 2, search.ArticleId)
	assert.Equal(t, "test", search.Content)
	assert.Equal(t, 1, search.Page)
	assert.Equal(t, 20, search.PageSize)
}

func TestCommentParise(t *testing.T) {
	parise := CommentParise{
		ID:     10,
		Parise: 1,
	}
	assert.Equal(t, 10, parise.ID)
	assert.Equal(t, 1, parise.Parise)
}
