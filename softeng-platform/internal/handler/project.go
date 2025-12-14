package handler

import (
	"net/http"
	"softeng-platform/internal/service"
	"softeng-platform/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService service.ProjectService
}

func NewProjectHandler(projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

// GetProjects 获取项目列表
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	category := c.Query("catagory")
	techStack := c.QueryArray("techStack")
	sort := c.Query("sort")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	cursor := c.Query("cursor")
	resourceType := c.Query("resourceType")

	projects, err := h.projectService.GetProjects(c.Request.Context(), category, techStack, sort, limit, cursor, resourceType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, projects)
}

// SearchProjects 搜索项目
func (h *ProjectHandler) SearchProjects(c *gin.Context) {
	keyword := c.Query("keyword")
	category := c.QueryArray("category")
	cursor := c.Query("cursor")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	projects, err := h.projectService.SearchProjects(c.Request.Context(), keyword, category, cursor, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, projects)
}

// GetProject 获取项目详情
func (h *ProjectHandler) GetProject(c *gin.Context) {
	projectID := c.Param("projectId")

	project, err := h.projectService.GetProject(c.Request.Context(), projectID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Project not found")
		return
	}

	response.Success(c, project)
}

// UpdateProject 更新项目
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	userID := c.GetInt("userID")
	projectID := c.Param("projectId")

	var req service.ProjectUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	result, err := h.projectService.UpdateProject(c.Request.Context(), userID, projectID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// UploadProject 上传项目
func (h *ProjectHandler) UploadProject(c *gin.Context) {
	userID := c.GetInt("userID")

	var req service.ProjectUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	result, err := h.projectService.UploadProject(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// LikeProject 点赞项目
func (h *ProjectHandler) LikeProject(c *gin.Context) {
	userID := c.GetInt("userID")
	projectID := c.Param("projectId")

	result, err := h.projectService.LikeProject(c.Request.Context(), userID, projectID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// UnlikeProject 取消点赞项目
func (h *ProjectHandler) UnlikeProject(c *gin.Context) {
	userID := c.GetInt("userID")
	projectID := c.Param("projectId")

	result, err := h.projectService.UnlikeProject(c.Request.Context(), userID, projectID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// AddComment 添加评论
func (h *ProjectHandler) AddComment(c *gin.Context) {
	userID := c.GetInt("userID")
	projectID := c.Param("projectId")

	var req struct {
		Content string `form:"content" json:"content" binding:"required"`
	}

	// 支持 multipart/form-data 和 application/json
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	result, err := h.projectService.AddComment(c.Request.Context(), userID, projectID, req.Content)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteComment 删除评论
func (h *ProjectHandler) DeleteComment(c *gin.Context) {
	userID := c.GetInt("userID")
	projectID := c.Param("projectId")

	result, err := h.projectService.DeleteComment(c.Request.Context(), userID, projectID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// ReplyComment 回复评论
func (h *ProjectHandler) ReplyComment(c *gin.Context) {
	userID := c.GetInt("userID")
	projectID := c.Param("projectId")
	commentID := c.Param("commentId")

	var req struct {
		Content string `form:"content" json:"content" binding:"required"`
	}

	// 支持 multipart/form-data 和 application/json
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	result, err := h.projectService.ReplyComment(c.Request.Context(), userID, projectID, commentID, req.Content)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteReply 删除回复
func (h *ProjectHandler) DeleteReply(c *gin.Context) {
	userID := c.GetInt("userID")
	projectID := c.Param("projectId")
	commentID := c.Param("commentId")

	result, err := h.projectService.DeleteReply(c.Request.Context(), userID, projectID, commentID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// AddView 增加项目浏览量
func (h *ProjectHandler) AddView(c *gin.Context) {
	projectID := c.Param("projectId")

	result, err := h.projectService.AddView(c.Request.Context(), projectID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// CollectProject 收藏项目
func (h *ProjectHandler) CollectProject(c *gin.Context) {
	userID := c.GetInt("userID")
	projectID := c.Param("projectId")

	result, err := h.projectService.CollectProject(c.Request.Context(), userID, projectID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// UncollectProject 取消收藏项目
func (h *ProjectHandler) UncollectProject(c *gin.Context) {
	userID := c.GetInt("userID")
	projectID := c.Param("projectId")

	result, err := h.projectService.UncollectProject(c.Request.Context(), userID, projectID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}
