package handler

import (
	"net/http"
	"softeng-platform/internal/service"
	"softeng-platform/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ToolHandler struct {
	toolService service.ToolService
}

func NewToolHandler(toolService service.ToolService) *ToolHandler {
	return &ToolHandler{toolService: toolService}
}

// GetTools 获取工具列表
func (h *ToolHandler) GetTools(c *gin.Context) {
	category := c.QueryArray("catagory")
	tags := c.QueryArray("tag")
	sort := c.Query("sort")
	cursor := c.Query("cursor")
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	tools, err := h.toolService.GetTools(c.Request.Context(), category, tags, sort, cursor, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, tools)
}

// SearchTools 搜索工具
func (h *ToolHandler) SearchTools(c *gin.Context) {
	keyword := c.Query("keyword")
	cursor := c.Query("cursor")
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	resourceType := c.Query("resourceType")

	tools, err := h.toolService.SearchTools(c.Request.Context(), keyword, cursor, pageSize, resourceType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, tools)
}

// GetTool 获取工具详情
func (h *ToolHandler) GetTool(c *gin.Context) {
	resourceID := c.Param("resourceId")
	resourceType := c.Query("resourceType")

	tool, err := h.toolService.GetTool(c.Request.Context(), resourceID, resourceType)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Tool not found")
		return
	}

	response.Success(c, tool)
}

// SubmitTool 提交工具
func (h *ToolHandler) SubmitTool(c *gin.Context) {
	userID := c.GetInt("userID")

	var req service.ToolSubmitRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	result, err := h.toolService.SubmitTool(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// LikeTool 点赞工具
func (h *ToolHandler) LikeTool(c *gin.Context) {
	userID := c.GetInt("userID")
	resourceID := c.Param("resourceId")

	result, err := h.toolService.LikeTool(c.Request.Context(), userID, resourceID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// UnlikeTool 取消点赞工具
func (h *ToolHandler) UnlikeTool(c *gin.Context) {
	userID := c.GetInt("userID")
	resourceID := c.Param("resourceId")

	result, err := h.toolService.UnlikeTool(c.Request.Context(), userID, resourceID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// CollectTool 收藏工具
func (h *ToolHandler) CollectTool(c *gin.Context) {
	userID := c.GetInt("userID")
	resourceID := c.Param("resourceId")
	resourceType := c.Query("resourceType")
	
	// resourceType 是必需的query参数
	if resourceType == "" {
		response.Error(c, http.StatusBadRequest, "resourceType is required")
		return
	}

	result, err := h.toolService.CollectTool(c.Request.Context(), userID, resourceID, resourceType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// UncollectTool 取消收藏工具
func (h *ToolHandler) UncollectTool(c *gin.Context) {
	userID := c.GetInt("userID")
	resourceID := c.Param("resourceId")
	resourceType := c.Query("resourceType")
	
	// resourceType 是必需的query参数
	if resourceType == "" {
		response.Error(c, http.StatusBadRequest, "resourceType is required")
		return
	}

	result, err := h.toolService.UncollectTool(c.Request.Context(), userID, resourceID, resourceType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// AddComment 添加评论
func (h *ToolHandler) AddComment(c *gin.Context) {
	userID := c.GetInt("userID")
	resourceID := c.Param("resourceId")
	resourceType := c.Query("resourceType")
	
	// resourceType 是必需的query参数
	if resourceType == "" {
		response.Error(c, http.StatusBadRequest, "resourceType is required")
		return
	}

	var req struct {
		Content string `form:"content" json:"content" binding:"required"`
	}

	// 支持 multipart/form-data 和 application/json
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	result, err := h.toolService.AddComment(c.Request.Context(), userID, resourceID, resourceType, req.Content)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteComment 删除评论
func (h *ToolHandler) DeleteComment(c *gin.Context) {
	userID := c.GetInt("userID")
	resourceID := c.Param("resourceId")

	result, err := h.toolService.DeleteComment(c.Request.Context(), userID, resourceID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// ReplyComment 回复评论
func (h *ToolHandler) ReplyComment(c *gin.Context) {
	userID := c.GetInt("userID")
	resourceID := c.Param("resourceId")
	commentID := c.Param("commentId")
	resourceType := c.Query("resourceType")
	
	// resourceType 是必需的query参数
	if resourceType == "" {
		response.Error(c, http.StatusBadRequest, "resourceType is required")
		return
	}

	var req struct {
		Content string `form:"content" json:"content" binding:"required"`
	}

	// 支持 multipart/form-data 和 application/json
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	result, err := h.toolService.ReplyComment(c.Request.Context(), userID, resourceID, commentID, resourceType, req.Content)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteReply 删除回复
func (h *ToolHandler) DeleteReply(c *gin.Context) {
	userID := c.GetInt("userID")
	resourceID := c.Param("resourceId")
	commentID := c.Param("commentId")

	result, err := h.toolService.DeleteReply(c.Request.Context(), userID, resourceID, commentID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// AddView 增加浏览量
func (h *ToolHandler) AddView(c *gin.Context) {
	resourceID := c.Param("resourceId")

	result, err := h.toolService.AddView(c.Request.Context(), resourceID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}
