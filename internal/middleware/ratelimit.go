package middleware

import (
	"net/http"
	"sync"
	"time"

	"softeng-platform/pkg/response"

	"github.com/gin-gonic/gin"
)

type slidingWindow struct {
	mu    sync.Mutex
	times []time.Time
}

// AuthRateLimit 限制每个客户端 IP 在固定时间窗口内的请求次数（用于 /auth 防刷）
func AuthRateLimit(max int, window time.Duration) gin.HandlerFunc {
	var m sync.Map
	return func(c *gin.Context) {
		if max <= 0 {
			c.Next()
			return
		}
		ip := c.ClientIP()
		v, _ := m.LoadOrStore(ip, &slidingWindow{})
		sw := v.(*slidingWindow)
		sw.mu.Lock()
		defer sw.mu.Unlock()
		now := time.Now()
		cutoff := now.Add(-window)
		var kept []time.Time
		for _, t := range sw.times {
			if t.After(cutoff) {
				kept = append(kept, t)
			}
		}
		sw.times = kept
		if len(sw.times) >= max {
			response.Error(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		sw.times = append(sw.times, now)
		c.Next()
	}
}
