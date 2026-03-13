package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// TestHelper 提供测试辅助功能
type TestHelper struct{}

// NewTestHelper 创建新的测试辅助实例
func NewTestHelper() *TestHelper {
	return &TestHelper{}
}

// CreateTestApp 创建用于测试的 Fiber 应用
func (h *TestHelper) CreateTestApp() *fiber.App {
	app := fiber.New()
	return app
}

// AssertResponseStatus 验证响应状态码
func (h *TestHelper) AssertResponseStatus(t *testing.T, resp *httptest.ResponseRecorder, expectedStatus int) {
	assert.Equal(t, expectedStatus, resp.Code)
}

// AssertResponseBody 验证响应体结构
func (h *TestHelper) AssertResponseBody(t *testing.T, resp *httptest.ResponseRecorder, expectedKeys ...string) {
	var responseBody map[string]any
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	for _, key := range expectedKeys {
		assert.Contains(t, responseBody, key)
	}
}

// AssertResponseMessage 验证响应消息
func (h *TestHelper) AssertResponseMessage(t *testing.T, resp *httptest.ResponseRecorder, expectedMessage string) {
	var responseBody map[string]any
	err := json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	if msg, exists := responseBody["msg"]; exists {
		assert.Equal(t, expectedMessage, msg)
	}
}

// CreateTestRequest 创建测试请求
func (h *TestHelper) CreateTestRequest(method, path string, body []byte) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

// MockHTTPClient 模拟 HTTP 客户端
type MockHTTPClient struct {
	Responses map[string]*MockResponse
}

type MockResponse struct {
	StatusCode int
	Body       []byte
	Error      error
}

func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{
		Responses: make(map[string]*MockResponse),
	}
}

func (m *MockHTTPClient) SetResponse(url string, response *MockResponse) {
	m.Responses[url] = response
}

// TestData 提供测试用的示例数据
type TestData struct{}

func NewTestData() *TestData {
	return &TestData{}
}

// GetSampleGithubData 获取示例 GitHub 数据
func (td *TestData) GetSampleGithubData() []map[string]any {
	return []map[string]any{
		{
			"author":     "testuser1",
			"message":    "feat: add new feature",
			"commitTime": "2024-01-01 12:00:00",
		},
		{
			"author":     "testuser2",
			"message":    "fix: resolve bug",
			"commitTime": "2024-01-02 13:00:00",
		},
	}
}

// GetSamplePageInfo 获取示例分页信息
func (td *TestData) GetSamplePageInfo() map[string]any {
	return map[string]any{
		"page":     1,
		"pageSize": 10,
	}
}
