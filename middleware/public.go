package middleware

import "strings"

// isPublicPath 公共资源路径，无需认证即可访问
func isPublicPath(path string) bool {
	return strings.Contains(path, "uploads/") ||
		strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") ||
		strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".gif") ||
		strings.HasSuffix(path, ".webp") || strings.HasSuffix(path, ".svg") ||
		strings.HasSuffix(path, ".ico") ||
		strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".js")
}
