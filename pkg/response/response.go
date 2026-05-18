package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse 错误响应格式，符合API文档规范
type ErrorResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"` // 可以为string或null
}

// Success 成功响应，data可以是包含message和其他字段的对象
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// Error 错误响应，符合API文档格式：{message: string, data: string|null}
func Error(c *gin.Context, code int, message string) {
	var data interface{} = nil
	// 如果message包含详细信息，可以将其放入data字段
	// 根据API文档，data可以是string或null
	c.JSON(code, ErrorResponse{
		Message: message,
		Data:    data,
	})
}

// ErrorWithData 带详细错误信息的错误响应
func ErrorWithData(c *gin.Context, code int, message string, errorData string) {
	c.JSON(code, ErrorResponse{
		Message: message,
		Data:    errorData,
	})
}
