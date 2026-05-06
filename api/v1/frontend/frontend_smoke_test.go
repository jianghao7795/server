package frontend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIHandlerStructsConstructible(t *testing.T) {
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
