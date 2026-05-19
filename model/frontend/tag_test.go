package frontend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrontendTagStruct(t *testing.T) {
	tag := Tag{
		Name:   "前端",
		Status: 1,
	}
	assert.Equal(t, "前端", tag.Name)
	assert.Equal(t, 1, tag.Status)
}

func TestFrontendTagTableName(t *testing.T) {
	assert.Equal(t, "tags", Tag{}.TableName())
}

func TestFrontendTagArticlesRelation(t *testing.T) {
	tag := Tag{
		Name: "React",
		Articles: []Article{
			{Title: "React基础"},
		},
	}
	assert.Len(t, tag.Articles, 1)
	assert.Equal(t, "React基础", tag.Articles[0].Title)
}
