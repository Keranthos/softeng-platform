package handler

import (
	"net/http"
	"softeng-platform/internal/model"
	"softeng-platform/internal/service"
	"softeng-platform/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	// 支持 multipart/form-data 和 application/json
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	token, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message":   "Registration successful",
		"JWT token": token,
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	// 支持 application/x-www-form-urlencoded 和 application/json
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message":   "1",
		"JWT token": token,
	})
}

// ForgotPassword 忘记密码
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email          string `form:"email" json:"email" binding:"required"`
		NewPassword    string `form:"new_password" json:"new_password" binding:"required"`
		CertifyPassword string `form:"certify_password" json:"certify_password" binding:"required"`
	}

	// 支持 multipart/form-data 和 application/json
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	err := h.authService.ForgotPassword(c.Request.Context(), req.Email, req.NewPassword, req.CertifyPassword)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Password reset successful",
	})
}
