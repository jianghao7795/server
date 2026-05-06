package ws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWsSingleton(t *testing.T) {
	assert.NotNil(t, WsSend)
}

func TestWsStructConstructible(t *testing.T) {
	var w Ws
	assert.NotNil(t, &w)
}
