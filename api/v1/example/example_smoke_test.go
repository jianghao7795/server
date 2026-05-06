package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIHandlerStructsConstructible(t *testing.T) {
	cases := []struct {
		name string
		ptr  any
	}{
		{"FileUploadAndDownloadApi", &FileUploadAndDownloadApi{}},
		{"ExcelApi", &ExcelApi{}},
		{"CustomerApi", &CustomerApi{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotNil(t, tc.ptr)
		})
	}
}
