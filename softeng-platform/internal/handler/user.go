package handler

import (
	"log"
	"net/http"
	"softeng-platform/internal/model"
	"softeng-platform/internal/service"
	"softeng-platform/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfile 获取个人资料
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt("userID")

	profile, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.Success(c, profile)
}

// Logout 用户登出
func (h *UserHandler) Logout(c *gin.Context) {
	// 在实际应用中，你可能需要将token加入黑名单
	// 这里简单返回成功消息
	response.Success(c, gin.H{
		"message": "Logout successful",
	})
}

// UpdateProfile 更新个人资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt("userID")

	var req model.UpdateProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

// GetCollection 获取个人收藏
func (h *UserHandler) GetCollection(c *gin.Context) {
	userID := c.GetInt("userID")

	collection, err := h.userService.GetCollection(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, collection)
}

// DeleteCollection 取消收藏
func (h *UserHandler) DeleteCollection(c *gin.Context) {
	userID := c.GetInt("userID")
	resourceType := c.Param("resourceType")
	resourceID := c.Param("resourceId")

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid resource ID")
		return
	}

	collection, err := h.userService.DeleteCollection(c.Request.Context(), userID, resourceType, resourceIDInt)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, collection)
}

// GetStatus 获取审核状态
func (h *UserHandler) GetStatus(c *gin.Context) {
	userID := c.GetInt("userID")

	status, err := h.userService.GetStatus(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, status)
}

// GetSummit 获取个人提交
func (h *UserHandler) GetSummit(c *gin.Context) {
	// 添加recover来捕获可能的panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[GetSummit] Panic recovered: %v", r)
			response.Error(c, http.StatusInternalServerError, "Internal server error")
		}
	}()

	userID := c.GetInt("userID")
	log.Printf("[GetSummit] 开始获取用户 %d 的提交记录", userID)

	summit, err := h.userService.GetSummit(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[GetSummit] 获取提交记录失败: %v", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 安全地获取长度，避免类型断言panic
	var toolsCount, coursesCount, projectsCount int
	if tools, ok := summit["tools"].([]map[string]interface{}); ok {
		toolsCount = len(tools)
	}
	if teaches, ok := summit["teaches"].([]map[string]interface{}); ok {
		coursesCount = len(teaches)
	}
	if resources, ok := summit["resources"].([]map[string]interface{}); ok {
		projectsCount = len(resources)
	}
	
	log.Printf("[GetSummit] 成功获取提交记录: tools=%d, courses=%d, projects=%d", 
		toolsCount, coursesCount, projectsCount)
	response.Success(c, summit)
}

// UpdateResourceStatus 更新资源状态
func (h *UserHandler) UpdateResourceStatus(c *gin.Context) {
	userID := c.GetInt("userID")
	resourceType := c.Param("resourceType")
	resourceID := c.Param("resourceId")

	var req struct {
		Action string `form:"action" json:"action" binding:"required"`
		State  string `form:"state" json:"state"`
	}

	// 支持 multipart/form-data 和 application/json
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	result, err := h.userService.UpdateResourceStatus(c.Request.Context(), userID, resourceType, resourceID, req.Action, req.State)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// UpdateEmail 更新邮箱
func (h *UserHandler) UpdateEmail(c *gin.Context) {
	userID := c.GetInt("userID")

	var req struct {
		Name     string `form:"name" json:"name" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
		NewEmail string `form:"new_email" json:"new_email" binding:"required"`
		Code     string `form:"code" json:"code" binding:"required"`
	}

	// 支持 multipart/form-data 和 application/json
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	user, err := h.userService.UpdateEmail(c.Request.Context(), userID, req.Name, req.Password, req.NewEmail, req.Code)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Email updated successfully",
		"user":    user,
	})
}

// UpdatePassword 更新密码
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID := c.GetInt("userID")

	var req struct {
		Name        string `form:"name" json:"name" binding:"required"`
		Email       string `form:"email" json:"email" binding:"required"`
		NewPassword string `form:"new_passward" json:"new_passward" binding:"required"`
		Code        string `form:"code" json:"code" binding:"required"`
	}

	// 支持 multipart/form-data 和 application/json
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	user, err := h.userService.UpdatePassword(c.Request.Context(), userID, req.Name, req.Email, req.NewPassword, req.Code)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Password updated successfully",
		"user":    user,
	})
}
