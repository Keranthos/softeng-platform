package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// UploadDir 上传文件目录
	UploadDir = "uploads/images"
	// MaxImageSize 最大图片大小（5MB）
	MaxImageSize = 5 * 1024 * 1024
	// AllowedImageTypes 允许的图片类型
	AllowedImageTypes = "image/jpeg,image/jpg,image/png,image/gif,image/webp"
)

// IsExternalURL 判断是否为外部URL
func IsExternalURL(url string) bool {
	if url == "" {
		return false
	}
	// 检查是否为Base64
	if strings.HasPrefix(url, "data:image/") {
		return false
	}
	// 检查是否为HTTP/HTTPS链接
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// IsBase64Image 判断是否为Base64图片
func IsBase64Image(url string) bool {
	return strings.HasPrefix(url, "data:image/")
}

// GenerateFileName 生成文件名（基于时间戳和MD5）
func GenerateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	if ext == "" {
		ext = ".jpg" // 默认扩展名
	}
	
	// 使用时间戳和原始名称生成MD5
	timestamp := time.Now().UnixNano()
	hash := md5.Sum([]byte(fmt.Sprintf("%d_%s", timestamp, originalName)))
	hashStr := hex.EncodeToString(hash[:])
	
	return fmt.Sprintf("%s%s", hashStr[:16], ext)
}

// GetUploadPath 获取上传路径（按年月组织）
func GetUploadPath() string {
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	return filepath.Join(UploadDir, year, month)
}

// EnsureUploadDir 确保上传目录存在
func EnsureUploadDir() error {
	path := GetUploadPath()
	return os.MkdirAll(path, 0755)
}

// SaveBase64Image 保存Base64图片到文件系统
func SaveBase64Image(base64Data string) (string, error) {
	// 解析Base64数据
	parts := strings.Split(base64Data, ",")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid base64 format")
	}
	
	// 获取图片类型
	header := parts[0]
	var ext string
	if strings.Contains(header, "jpeg") || strings.Contains(header, "jpg") {
		ext = ".jpg"
	} else if strings.Contains(header, "png") {
		ext = ".png"
	} else if strings.Contains(header, "gif") {
		ext = ".gif"
	} else if strings.Contains(header, "webp") {
		ext = ".webp"
	} else {
		ext = ".jpg" // 默认
	}
	
	// 解码Base64
	data := parts[1]
	imageData, err := decodeBase64(data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}
	
	// 确保目录存在
	if err := EnsureUploadDir(); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}
	
	// 生成文件名
	fileName := GenerateFileName("base64" + ext)
	uploadPath := GetUploadPath()
	filePath := filepath.Join(uploadPath, fileName)
	
	// 保存文件
	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}
	
	// 返回相对路径（用于URL）
	relativePath := filepath.Join(uploadPath, fileName)
	// 将路径分隔符统一为 /
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	return "/" + relativePath, nil
}

// DownloadAndSaveImage 下载外部图片并保存到本地
func DownloadAndSaveImage(url string) (string, error) {
	// 创建HTTP客户端，设置超时
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// 下载图片
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()
	
	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
	}
	
	// 检查Content-Type
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("invalid content type: %s", contentType)
	}
	
	// 检查文件大小
	if resp.ContentLength > MaxImageSize {
		return "", fmt.Errorf("image too large: %d bytes", resp.ContentLength)
	}
	
	// 读取图片数据
	imageData, err := io.ReadAll(io.LimitReader(resp.Body, MaxImageSize+1))
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}
	
	// 再次检查大小（防止Content-Length不准确）
	if len(imageData) > MaxImageSize {
		return "", fmt.Errorf("image too large: %d bytes", len(imageData))
	}
	
	// 从URL获取扩展名
	ext := filepath.Ext(url)
	if ext == "" {
		// 从Content-Type推断扩展名
		if strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg") {
			ext = ".jpg"
		} else if strings.Contains(contentType, "png") {
			ext = ".png"
		} else if strings.Contains(contentType, "gif") {
			ext = ".gif"
		} else if strings.Contains(contentType, "webp") {
			ext = ".webp"
		} else {
			ext = ".jpg" // 默认
		}
	}
	
	// 确保目录存在
	if err := EnsureUploadDir(); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}
	
	// 生成文件名
	fileName := GenerateFileName("external" + ext)
	uploadPath := GetUploadPath()
	filePath := filepath.Join(uploadPath, fileName)
	
	// 保存文件
	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}
	
	// 返回相对路径（用于URL）
	relativePath := filepath.Join(uploadPath, fileName)
	// 将路径分隔符统一为 /
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	return "/" + relativePath, nil
}

// ProcessImageURL 处理图片URL（自动本地化外部链接和Base64）
func ProcessImageURL(imageURL string) (string, error) {
	if imageURL == "" {
		return "", nil
	}
	
	// 如果是Base64，保存到本地
	if IsBase64Image(imageURL) {
		localPath, err := SaveBase64Image(imageURL)
		if err != nil {
			return "", fmt.Errorf("failed to save base64 image: %w", err)
		}
		return localPath, nil
	}
	
	// 如果是外部URL，下载到本地
	if IsExternalURL(imageURL) {
		localPath, err := DownloadAndSaveImage(imageURL)
		if err != nil {
			// 如果下载失败，返回原URL（不强制本地化）
			return imageURL, nil
		}
		return localPath, nil
	}
	
	// 如果已经是本地路径，直接返回
	return imageURL, nil
}

// DeleteImageFile 删除图片文件
func DeleteImageFile(imageURL string) error {
	if imageURL == "" {
		return nil
	}
	
	// 只删除本地路径的图片（以 /uploads/ 开头）
	if !strings.HasPrefix(imageURL, "/uploads/") {
		// 外部URL或Base64，不删除
		return nil
	}
	
	// 移除开头的 /，转换为文件系统路径
	filePath := strings.TrimPrefix(imageURL, "/")
	// 将路径分隔符统一为系统分隔符
	filePath = filepath.FromSlash(filePath)
	
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 文件不存在，忽略错误
		return nil
	}
	
	// 删除文件
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete image file: %w", err)
	}
	
	return nil
}

// decodeBase64 解码Base64字符串
func decodeBase64(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

