package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExampleAPIHandlerStructsConstructible(t *testing.T) {
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

func TestExampleServiceVars(t *testing.T) {
	assert.NotNil(t, excelService)
	assert.NotNil(t, customerService)
	assert.NotNil(t, fileUploadAndDownloadService)
}
