package handler

import (
	"net/http"
	"path/filepath"
	"softeng-platform/internal/utils"
	"softeng-platform/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// UploadImage 上传图片
// @Summary 上传图片
// @Description 支持上传图片文件，返回本地URL
// @Tags 上传
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "图片文件"
// @Success 200 {object} map[string]interface{} "成功返回图片URL"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/upload/image [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("image")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择图片文件")
		return
	}

	// 检查文件大小
	if file.Size > utils.MaxImageSize {
		response.Error(c, http.StatusBadRequest, "图片大小不能超过5MB")
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	isAllowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		response.Error(c, http.StatusBadRequest, "不支持的图片格式，仅支持 JPG、PNG、GIF、WebP")
		return
	}

	// 确保上传目录存在
	if err := utils.EnsureUploadDir(); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建上传目录失败")
		return
	}

	// 生成文件名
	fileName := utils.GenerateFileName(file.Filename)
	uploadPath := utils.GetUploadPath()
	filePath := filepath.Join(uploadPath, fileName)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		response.Error(c, http.StatusInternalServerError, "保存文件失败")
		return
	}

	// 返回相对路径（用于URL）
	relativePath := filepath.Join(uploadPath, fileName)
	// 将路径分隔符统一为 /
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	imageURL := "/" + relativePath

	response.Success(c, map[string]interface{}{
		"url": imageURL,
	})
}

// ProcessImageURL 处理图片URL（自动本地化）
// @Summary 处理图片URL
// @Description 自动将外部URL或Base64图片本地化
// @Tags 上传
// @Accept json
// @Produce json
// @Param request body map[string]string true "图片URL"
// @Success 200 {object} map[string]interface{} "成功返回本地URL"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/upload/process [post]
func (h *UploadHandler) ProcessImageURL(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请提供图片URL")
		return
	}

	// 处理图片URL（自动本地化）
	localURL, err := utils.ProcessImageURL(req.URL)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "处理图片失败: "+err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"url": localURL,
	})
}

