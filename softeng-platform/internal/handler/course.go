package handler

import (
	"net/http"
	"softeng-platform/internal/service"
	"softeng-platform/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	courseService service.CourseService
}

func NewCourseHandler(courseService service.CourseService) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

// GetCourses 获取课程列表
func (h *CourseHandler) GetCourses(c *gin.Context) {
	semester := c.Query("semester")
	category := c.QueryArray("category")
	sort := c.Query("sort")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	cursor, _ := strconv.Atoi(c.Query("cursor"))
	resourceType := c.Query("resourceType")

	courses, err := h.courseService.GetCourses(c.Request.Context(), semester, category, sort, limit, cursor, resourceType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, courses)
}

// SearchCourses 搜索课程
func (h *CourseHandler) SearchCourses(c *gin.Context) {
	keyword := c.Query("keyword")
	category := c.QueryArray("category")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	cursor, _ := strconv.Atoi(c.Query("cursor"))
	resourceType := c.Query("resourceType")

	courses, err := h.courseService.SearchCourses(c.Request.Context(), keyword, category, limit, cursor, resourceType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, courses)
}

// GetCourse 获取课程详情
func (h *CourseHandler) GetCourse(c *gin.Context) {
	courseID := c.Param("courseId")
	resourceType := c.Query("resourceType")

	course, err := h.courseService.GetCourse(c.Request.Context(), courseID, resourceType)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Course not found")
		return
	}

	response.Success(c, course)
}

// UploadResource 上传课程资源
func (h *CourseHandler) UploadResource(c *gin.Context) {
	userID := c.GetInt("userID")
	courseID := c.Param("courseId")
	resourceType := c.Query("resourceType")
	
	// resourceType 是必需的query参数
	if resourceType == "" {
		response.Error(c, http.StatusBadRequest, "resourceType is required")
		return
	}

	var req service.CourseUploadRequest
	// 支持 multipart/form-data 和 application/json
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	result, err := h.courseService.UploadResource(c.Request.Context(), userID, courseID, resourceType, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// DownloadTextbook 下载课本
func (h *CourseHandler) DownloadTextbook(c *gin.Context) {
	courseID := c.Param("courseId")
	textbookID := c.Param("textbookId")

	result, err := h.courseService.DownloadTextbook(c.Request.Context(), courseID, textbookID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// AddComment 发表评论
func (h *CourseHandler) AddComment(c *gin.Context) {
	userID := c.GetInt("userID")
	courseID := c.Param("courseId")

	var req struct {
		Content string `form:"content" json:"content" binding:"required"`
	}

	// 支持 multipart/form-data 和 application/json
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	result, err := h.courseService.AddComment(c.Request.Context(), userID, courseID, req.Content)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteComment 删除评论
func (h *CourseHandler) DeleteComment(c *gin.Context) {
	userID := c.GetInt("userID")
	courseID := c.Param("courseId")

	result, err := h.courseService.DeleteComment(c.Request.Context(), userID, courseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// ReplyComment 回复评论
func (h *CourseHandler) ReplyComment(c *gin.Context) {
	userID := c.GetInt("userID")
	courseID := c.Param("courseId")
	commentID := c.Param("commentId")

	var req struct {
		Content string `form:"content" json:"content" binding:"required"`
	}

	// 支持 multipart/form-data 和 application/json
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	result, err := h.courseService.ReplyComment(c.Request.Context(), userID, courseID, commentID, req.Content)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteReply 删除回复
func (h *CourseHandler) DeleteReply(c *gin.Context) {
	userID := c.GetInt("userID")
	courseID := c.Param("courseId")
	commentID := c.Param("commentId")

	result, err := h.courseService.DeleteReply(c.Request.Context(), userID, courseID, commentID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// AddView 增加课程浏览量
func (h *CourseHandler) AddView(c *gin.Context) {
	courseID := c.Param("courseId")

	result, err := h.courseService.AddView(c.Request.Context(), courseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// CollectCourse 收藏课程
func (h *CourseHandler) CollectCourse(c *gin.Context) {
	userID := c.GetInt("userID")
	courseID := c.Param("courseId")

	result, err := h.courseService.CollectCourse(c.Request.Context(), userID, courseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// UncollectCourse 取消收藏课程
func (h *CourseHandler) UncollectCourse(c *gin.Context) {
	userID := c.GetInt("userID")
	courseID := c.Param("courseId")

	result, err := h.courseService.UncollectCourse(c.Request.Context(), userID, courseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// LikeCourse 点赞课程
func (h *CourseHandler) LikeCourse(c *gin.Context) {
	userID := c.GetInt("userID")
	courseID := c.Param("courseId")

	result, err := h.courseService.LikeCourse(c.Request.Context(), userID, courseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// UnlikeCourse 取消点赞课程
func (h *CourseHandler) UnlikeCourse(c *gin.Context) {
	userID := c.GetInt("userID")
	courseID := c.Param("courseId")

	result, err := h.courseService.UnlikeCourse(c.Request.Context(), userID, courseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}
