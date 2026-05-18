package handler

import (
	"net/http"
	"softeng-platform/internal/service"
	"softeng-platform/pkg/response"

	"github.com/gin-gonic/gin"
)

type RagAgentHandler struct {
	svc service.RagAgentService
}

func NewRagAgentHandler(svc service.RagAgentService) *RagAgentHandler {
	return &RagAgentHandler{svc: svc}
}

type ragAgentChatRequest struct {
	Message string `json:"message"`
}

// Chat POST /api/agent/rag — 检索平台工具/课程/项目后由 LLM 或离线模板生成回答（演示用，无需登录）
func (h *RagAgentHandler) Chat(c *gin.Context) {
	var req ragAgentChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求体需为 JSON，且包含 message 字段")
		return
	}
	res, err := h.svc.Chat(c.Request.Context(), req.Message)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, res)
}
