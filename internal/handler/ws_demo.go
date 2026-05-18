package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"softeng-platform/internal/utils"
	"softeng-platform/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true },
}

// WsDemoStream GET /api/ws/stream — 演示用 WebSocket（JSON 心跳）
// 安全：须 query ?token=JWT 且用户为 admin/superadmin（浏览器 WebSocket 无法设 Header，故用 query）
func WsDemoStream(c *gin.Context) {
	token := strings.TrimSpace(c.Query("token"))
	if token == "" {
		response.Error(c, http.StatusUnauthorized, "token query required")
		return
	}
	claims, err := utils.ValidateToken(token)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid token")
		return
	}
	if claims.Role != "admin" && claims.Role != "superadmin" {
		response.Error(c, http.StatusForbidden, "admin access required")
		return
	}

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				cancel()
				return
			}
		}
	}()

	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	send := func() {
		payload := map[string]string{
			"type":  "demo",
			"title": "服务端 WebSocket",
			"body":  "心跳 " + time.Now().Format("15:04:05") + "（可扩展为审核结果推送）",
		}
		b, _ := json.Marshal(payload)
		_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		_ = conn.WriteMessage(websocket.TextMessage, b)
	}
	send()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			send()
		}
	}
}
