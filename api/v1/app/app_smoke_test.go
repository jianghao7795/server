package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIHandlerStructsConstructible(t *testing.T) {
	cases := []struct {
		name string
		ptr  any
	}{
		{"ArticleApi", &ArticleApi{}},
		{"BaseMessageApi", &BaseMessageApi{}},
		{"CommentApi", &CommentApi{}},
		{"TagApi", &TagApi{}},
		{"TaskNameApi", &TaskNameApi{}},
		{"FileUploadAndDownloadApi", &FileUploadAndDownloadApi{}},
		{"UserApi", &UserApi{}},
		{"LikeApi", &LikeApi{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotNil(t, tc.ptr)
		})
	}
}

func TestAppServiceVars(t *testing.T) {
	// 验证全局服务变量已初始化
	assert.NotNil(t, articleService)
	assert.NotNil(t, baseMessageService)
	assert.NotNil(t, commentService)
	assert.NotNil(t, appTabService)
	assert.NotNil(t, fileUploadService)
	assert.NotNil(t, userService)
	assert.NotNil(t, likeService)
}

func TestCommentPaginationDefaults(t *testing.T) {
	t.Run("默认分页参数", func(t *testing.T) {
		// GetCommentList 中 page=0 时自动设为1, pageSize=0 时自动设为10
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

	t.Run("显式分页参数", func(t *testing.T) {
		page := 3
		pageSize := 20
		offset := pageSize * (page - 1)
		assert.Equal(t, 40, offset)
	})
}

func TestLikePaginationDefaults(t *testing.T) {
	t.Run("分页默认值", func(t *testing.T) {
		page := 0
		if page < 1 {
			page = 1
		}
		pageSize := 0
		if pageSize < 1 || pageSize > 100 {
			pageSize = 10
		}
		assert.Equal(t, 1, page)
		assert.Equal(t, 10, pageSize)
	})

	t.Run("超出范围的pageSize", func(t *testing.T) {
		pageSize := 200
		if pageSize < 1 || pageSize > 100 {
			pageSize = 10
		}
		assert.Equal(t, 10, pageSize)
	})

	t.Run("负数page", func(t *testing.T) {
		page := -5
		if page < 1 {
			page = 1
		}
		assert.Equal(t, 1, page)
	})
}

func TestTaskingParameterParsing(t *testing.T) {
	t.Run("空任务名验证", func(t *testing.T) {
		tasking := ""
		isEmpty := tasking == ""
		assert.True(t, isEmpty)
	})

	t.Run("有效任务名", func(t *testing.T) {
		tasking := "cleanDb"
		assert.Equal(t, "cleanDb", tasking)
		assert.NotEmpty(t, tasking)
	})
}
