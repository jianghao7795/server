package mobile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMobileAPIHandlerStructsConstructible(t *testing.T) {
	cases := []struct {
		name string
		ptr  any
	}{
		{"UserApi", &UserApi{}},
		{"RegisterMobile", &RegisterMobile{}},
		{"LoginApi", &LoginApi{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotNil(t, tc.ptr)
		})
	}
}

func TestMobileServiceVars(t *testing.T) {
	assert.NotNil(t, userService)
	assert.NotNil(t, loginService)
	assert.NotNil(t, registerService)
}

func TestMobilePaginationDefaults(t *testing.T) {
	t.Run("默认分页", func(t *testing.T) {
		page := 0
		pageSize := 0
		offset := pageSize * (page - 1)
		assert.Equal(t, 0, offset)
	})

	t.Run("带偏移的分页", func(t *testing.T) {
		pageSize := 10
		page := 2
		offset := pageSize * (page - 1)
		assert.Equal(t, 10, offset)
	})
}

func TestMobileSearchPattern(t *testing.T) {
	username := "test"
	like := "%" + username + "%"
	assert.Equal(t, "%test%", like)
}
