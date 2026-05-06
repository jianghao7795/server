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
