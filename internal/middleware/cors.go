package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS 按允许的 Origin 白名单设置跨域头；白名单来自环境变量 CORS_ALLOWED_ORIGINS（逗号分隔）
func CORS(allowedOrigins []string) gin.HandlerFunc {
	allowed := allowedOrigins
	if len(allowed) == 0 {
		allowed = []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://localhost:3001",
			"http://127.0.0.1:3001",
		}
	}

	return func(c *gin.Context) {
		origin := strings.TrimSpace(c.Request.Header.Get("Origin"))

		if origin != "" {
			for _, o := range allowed {
				if origin == strings.TrimSpace(o) {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
