package request

import (
	"testing"

	"server/model/app"

	"github.com/stretchr/testify/assert"
)

func TestTagSearchStruct(t *testing.T) {
	search := TagSearch{
		Tag: app.Tag{
			Name:   "Go",
			Status: 1,
		},
	}
	search.Page = 1
	search.PageSize = 10

	assert.Equal(t, "Go", search.Name)
	assert.Equal(t, 1, search.Status)
	assert.Equal(t, 1, search.Page)
	assert.Equal(t, 10, search.PageSize)
}

func TestTagSearchDefaultValues(t *testing.T) {
	search := TagSearch{}
	assert.Empty(t, search.Name)
	assert.Zero(t, search.Status)
	assert.Zero(t, search.Page)
	assert.Zero(t, search.PageSize)
}
