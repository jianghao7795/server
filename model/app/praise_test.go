package app

import (
	"testing"

	global "server/model"

	"github.com/stretchr/testify/assert"
)

func TestPraiseStruct(t *testing.T) {
	praise := global.Praise{
		CommentId: 1,
		UserId:    100,
	}
	assert.Equal(t, uint(1), praise.CommentId)
	assert.Equal(t, uint(100), praise.UserId)
}

func TestPraiseTableName(t *testing.T) {
	assert.Equal(t, "praise", global.Praise{}.TableName())
}
