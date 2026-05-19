package request

import (
	"testing"

	"server/model/app"

	"github.com/stretchr/testify/assert"
)

func TestFrontendTagSearchStruct(t *testing.T) {
	search := TagSearch{
		Tag: app.Tag{
			Name:   "React",
			Status: 1,
		},
	}
	search.Page = 2
	search.PageSize = 15

	assert.Equal(t, "React", search.Name)
	assert.Equal(t, 1, search.Status)
	assert.Equal(t, 2, search.Page)
	assert.Equal(t, 15, search.PageSize)
}
