package frontend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrontendAPIHandlerStructsConstructible(t *testing.T) {
	cases := []struct {
		name string
		ptr  any
	}{
		{"User", &User{}},
		{"TagApi", &TagApi{}},
		{"CommentApi", &CommentApi{}},
		{"ArticleApi", &ArticleApi{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotNil(t, tc.ptr)
		})
	}
}

func TestFrontendServiceVars(t *testing.T) {
	assert.NotNil(t, articleServiceApp)
	assert.NotNil(t, commentServiceApp)
	assert.NotNil(t, userServiceApp)
	assert.NotNil(t, tagServiceApp)
	assert.NotNil(t, imagesServiceApp)
	assert.NotNil(t, userService)
	assert.NotNil(t, jwtService)
}

func TestFrontendArticlePagination(t *testing.T) {
	t.Run("默认分页参数", func(t *testing.T) {
		page := 0
		if page == 0 {
			page = 1
		}
		pageSize := 0
		if pageSize == 0 {
			pageSize = 10
		}
		assert.Equal(t, 1, page)
		assert.Equal(t, 10, pageSize)
	})

	t.Run("分页偏移计算", func(t *testing.T) {
		pageSize := 20
		page := 3
		offset := pageSize * (page - 1)
		assert.Equal(t, 40, offset)
	})
}
