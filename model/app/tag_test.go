package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagStruct(t *testing.T) {
	tag := Tag{
		Name:   "技术",
		Status: 1,
	}
	assert.Equal(t, "技术", tag.Name)
	assert.Equal(t, 1, tag.Status)
}

func TestTagTableName(t *testing.T) {
	assert.Equal(t, "tags", Tag{}.TableName())
}

func TestTagArticlesRelation(t *testing.T) {
	tag := Tag{
		Name: "Go",
		Articles: []Article{
			{Title: "Go入门"},
			{Title: "Go进阶"},
		},
	}
	assert.Len(t, tag.Articles, 2)
}
