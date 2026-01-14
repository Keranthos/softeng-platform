package handler

import (
	"fmt"
	"log"
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
// @Summary 获取工具列表
// @Description 根据分类、标签等条件获取工具列表，支持分页和排序
// @Tags tools
// @Accept json
// @Produce json
// @Param catagory query array false "工具分类"
// @Param tag query array false "工具标签"
// @Param sort query string false "排序方式：最新/newest, 最多浏览/views, 最多收藏/collections, 最多点赞/loves"
// @Param cursor query string false "分页游标"
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{} "成功返回工具列表"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /tools [get]
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
// @Summary 搜索工具
// @Description 根据关键词搜索工具
// @Tags tools
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param cursor query string false "分页游标"
// @Param page_size query int false "每页数量" default(10)
// @Param resourceType query string false "资源类型"
// @Success 200 {object} map[string]interface{} "成功返回搜索结果"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /tools/search [get]
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
// @Summary 获取工具详情
// @Description 根据工具ID获取详细信息，包括评论、标签等
// @Tags tools
// @Accept json
// @Produce json
// @Param resourceId path string true "工具ID"
// @Param resourceType query string false "资源类型"
// @Success 200 {object} map[string]interface{} "成功返回工具详情"
// @Failure 404 {object} map[string]interface{} "工具不存在"
// @Router /tools/{resourceId} [get]
func (h *ToolHandler) GetTool(c *gin.Context) {
	resourceID := c.Param("resourceId")
	resourceType := c.Query("resourceType")
	
	log.Printf("[DEBUG] GetTool called: resourceID=%s, resourceType=%s", resourceID, resourceType)
	
	// 获取可选的用户ID（如果用户已登录）
	userID := 0
	if uid, exists := c.Get("userID"); exists {
		if id, ok := uid.(int); ok {
			userID = id
		}
	}

	tool, err := h.toolService.GetTool(c.Request.Context(), resourceID, resourceType, userID)
	if err != nil {
		// 记录详细错误信息用于调试
		log.Printf("[DEBUG] GetTool error: resourceID=%s, error=%v", resourceID, err)
		response.Error(c, http.StatusNotFound, fmt.Sprintf("Tool not found: %v", err))
		return
	}

	log.Printf("[DEBUG] GetTool success: resourceID=%s", resourceID)
	response.Success(c, tool)
}

// SubmitTool 提交工具
func (h *ToolHandler) SubmitTool(c *gin.Context) {
	userID := c.GetInt("userID")

	var req service.ToolSubmitRequest
	// 根据 Content-Type 选择绑定方式
	contentType := c.GetHeader("Content-Type")
	
	var bindErr error
	if contentType == "application/json" {
		bindErr = c.ShouldBindJSON(&req)
	} else {
		// 对于 form 数据，需要手动处理数组字段
		// 先绑定非数组字段
		req.Name = c.PostForm("name")
		req.Link = c.PostForm("link")
		req.Description = c.PostForm("description")
		req.DescriptionDetail = c.PostForm("description_detail")
		req.Category = c.PostForm("catagory")
		req.ToolType = c.PostForm("type")
		
		// 手动获取数组字段
		req.Tags = c.PostFormArray("tags")
		req.Images = c.PostFormArray("images")
		
		// 验证必填字段
		if req.Name == "" || req.Link == "" || req.Description == "" || req.DescriptionDetail == "" || req.Category == "" || len(req.Tags) == 0 {
			log.Printf("[DEBUG] SubmitTool validation failed: Name=%s, Link=%s, Description=%s, DescriptionDetail=%s, Category=%s, Tags=%v", 
				req.Name, req.Link, req.Description, req.DescriptionDetail, req.Category, req.Tags)
			response.Error(c, http.StatusBadRequest, "Missing required fields")
			return
		}
	}
	
	if bindErr != nil {
		log.Printf("[DEBUG] SubmitTool binding error: %v", bindErr)
		log.Printf("[DEBUG] SubmitTool Content-Type: %s", contentType)
		log.Printf("[DEBUG] SubmitTool userID: %d", userID)
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid request data: %v", bindErr))
		return
	}

	// 验证必填字段（JSON 请求）
	if contentType == "application/json" {
		if req.Name == "" || req.Link == "" || req.Description == "" || req.DescriptionDetail == "" || req.Category == "" || len(req.Tags) == 0 {
			log.Printf("[DEBUG] SubmitTool validation failed: Name=%s, Link=%s, Description=%s, DescriptionDetail=%s, Category=%s, Tags=%v", 
				req.Name, req.Link, req.Description, req.DescriptionDetail, req.Category, req.Tags)
			response.Error(c, http.StatusBadRequest, "Missing required fields")
			return
		}
	}

	log.Printf("[DEBUG] SubmitTool success: Name=%s, Tags=%v, Images=%v", req.Name, req.Tags, req.Images)

	result, err := h.toolService.SubmitTool(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// UpdateTool 更新工具
func (h *ToolHandler) UpdateTool(c *gin.Context) {
	userID := c.GetInt("userID")
	resourceID := c.Param("resourceId")

	var req service.ToolSubmitRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	result, err := h.toolService.UpdateTool(c.Request.Context(), userID, resourceID, req)
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
	
	// resourceType可以从query参数或POST body中获取，如果没有则默认为'tool'
	resourceType := c.Query("resourceType")
	if resourceType == "" {
		resourceType = c.PostForm("resourceType")
	}
	if resourceType == "" {
		resourceType = "tool" // 默认值，因为这是tools路由
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
	
	// resourceType可以从query参数中获取，如果没有则默认为'tool'
	// 注意：DELETE请求的body在Gin中较难获取，优先使用query参数
	resourceType := c.Query("resourceType")
	if resourceType == "" {
		// 尝试从request body中获取（某些客户端可能在DELETE body中传递参数）
		if c.Request.ContentLength > 0 {
			var body map[string]interface{}
			if err := c.ShouldBindJSON(&body); err == nil {
				if rt, ok := body["resourceType"].(string); ok {
					resourceType = rt
				}
			}
		}
	}
	if resourceType == "" {
		resourceType = "tool" // 默认值
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
	
	// resourceType可以从query参数或POST body中获取，如果没有则默认为'tool'
	resourceType := c.Query("resourceType")
	if resourceType == "" {
		resourceType = c.PostForm("resourceType")
	}
	if resourceType == "" {
		resourceType = "tool" // 默认值
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
	commentID := c.Param("commentId")

	result, err := h.toolService.DeleteComment(c.Request.Context(), userID, resourceID, commentID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// GetComments 获取工具评论列表
func (h *ToolHandler) GetComments(c *gin.Context) {
	resourceID := c.Param("resourceId")
	cursor := c.Query("cursor")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	comments, err := h.toolService.GetComments(c.Request.Context(), resourceID, cursor, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, comments)
}

// LikeComment 点赞/取消点赞评论
func (h *ToolHandler) LikeComment(c *gin.Context) {
	userID := c.GetInt("userID")
	resourceID := c.Param("resourceId")
	commentID := c.Param("commentId")

	result, err := h.toolService.LikeComment(c.Request.Context(), userID, resourceID, commentID)
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
	
	// resourceType可以从query参数或POST body中获取，如果没有则默认为'tool'
	resourceType := c.Query("resourceType")
	if resourceType == "" {
		resourceType = c.PostForm("resourceType")
	}
	if resourceType == "" {
		resourceType = "tool" // 默认值
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
