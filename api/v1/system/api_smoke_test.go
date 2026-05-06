package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIHandlerStructsConstructible(t *testing.T) {
	cases := []struct {
		name string
		ptr  any
	}{
		{"BaseApi", &BaseApi{}},
		{"SystemApi", &SystemApi{}},
		{"SystemGithubApi", &SystemGithubApi{}},
		{"OperationRecordApi", &OperationRecordApi{}},
		{"AuthorityMenuApi", &AuthorityMenuApi{}},
		{"JwtApi", &JwtApi{}},
		{"DBApi", &DBApi{}},
		{"DictionaryApi", &DictionaryApi{}},
		{"DictionaryDetailApi", &DictionaryDetailApi{}},
		{"CasbinApi", &CasbinApi{}},
		{"AutoCodeApi", &AutoCodeApi{}},
		{"AutoCodeHistoryApi", &AutoCodeHistoryApi{}},
		{"AuthorityApi", &AuthorityApi{}},
		{"AuthorityBtnApi", &AuthorityBtnApi{}},
		{"SystemApiApi", &SystemApiApi{}},
		{"UserProblem", &UserProblem{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotNil(t, tc.ptr)
		})
	}
}
