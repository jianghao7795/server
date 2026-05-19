package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPraiseStruct(t *testing.T) {
	praise := Praise{
		CommentId: 1,
		UserId:    100,
	}
	assert.Equal(t, int64(1), praise.CommentId)
	assert.Equal(t, int64(100), praise.UserId)
}

func TestPraiseTableName(t *testing.T) {
	assert.Equal(t, "praise", Praise{}.TableName())
}
