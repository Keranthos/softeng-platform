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
	// 支持前端传递的 page 和 page_size 参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	// 如果没有page参数，尝试从cursor获取
	cursor, _ := strconv.Atoi(c.Query("cursor"))
	if cursor == 0 {
		cursor = (page - 1) * pageSize
	}
	
	limit := pageSize
	if limit == 0 {
		limit = 10
	}
	
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
		Action       string `form:"action" json:"action" query:"action"`
		ResourceType string `form:"resourceType" json:"resourceType" query:"resourceType"`
		RejectReason string `form:"reject_reason" json:"reject_reason" query:"reject_reason"`
	}

	// 支持 GET (query参数) 和 POST (form/json)
	if c.Request.Method == "GET" {
		req.Action = c.Query("action")
		req.ResourceType = c.Query("resourceType")
		req.RejectReason = c.Query("reject_reason")
	} else {
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid request data")
			return
		}
	}

	if req.Action == "" {
		response.Error(c, http.StatusBadRequest, "action is required")
		return
	}

	// 转换前端action值：approve -> approved, reject -> rejected
	status := req.Action
	if req.Action == "approve" {
		status = "approved"
	} else if req.Action == "reject" {
		status = "rejected"
	}

	err := h.adminService.ReviewItem(c.Request.Context(), itemID, status, req.ResourceType, req.RejectReason)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Review completed successfully",
	})
}
