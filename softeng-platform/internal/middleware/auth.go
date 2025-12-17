package middleware

import (
	"net/http"
	"softeng-platform/internal/utils"
	"softeng-platform/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			response.Error(c, http.StatusForbidden, "Admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}
