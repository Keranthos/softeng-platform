package handler

import (
	"net/http"
	"softeng-platform/internal/service"
	"softeng-platform/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// GetPending 获取待审核内容
func (h *AdminHandler) GetPending(c *gin.Context) {
	itemType := c.Query("type")
	cursor, _ := strconv.Atoi(c.Query("cursor"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sort := c.Query("sort")

	result, err := h.adminService.GetPending(c.Request.Context(), itemType, cursor, limit, sort)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// ReviewItem 审核项目
func (h *AdminHandler) ReviewItem(c *gin.Context) {
	itemID := c.Param("itemId")

	var req struct {
		Action       string `form:"action" json:"action" binding:"required"`
		RejectReason string `form:"gejrct_reason" json:"gejrct_reason"` // 保持与API文档一致的拼写
	}

	// 支持 multipart/form-data 和 application/json
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	err := h.adminService.ReviewItem(c.Request.Context(), itemID, req.Action, req.RejectReason)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Review completed successfully",
	})
}
