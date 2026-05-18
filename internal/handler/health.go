package handler

import (
	"context"
	"net/http"
	"time"

	"softeng-platform/internal/repository"

	"github.com/gin-gonic/gin"
)

// Live 进程存活探针（不访问数据库，供编排做 liveness）
func Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Health 就绪探针：检查 MySQL 连通（短超时，供编排 readiness / 负载均衡）
func Health(db *repository.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "unhealthy",
				"database": "down",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": "up",
		})
	}
}
