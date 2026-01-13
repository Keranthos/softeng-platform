package service

import (
	"context"
	"fmt"
	"softeng-platform/internal/repository"
	"softeng-platform/internal/utils"
)

type ToolService interface {
	GetTools(ctx context.Context, category, tags []string, sort, cursor string, pageSize int) (map[string]interface{}, error)
	GetTool(ctx context.Context, resourceID, resourceType string, userID int) (map[string]interface{}, error)
	SearchTools(ctx context.Context, keyword, cursor string, pageSize int, resourceType string) (map[string]interface{}, error)
	SubmitTool(ctx context.Context, userID int, req ToolSubmitRequest) (map[string]interface{}, error)
	LikeTool(ctx context.Context, userID int, resourceID string) (map[string]interface{}, error)
	UnlikeTool(ctx context.Context, userID int, resourceID string) (map[string]interface{}, error)
	CollectTool(ctx context.Context, userID int, resourceID, resourceType string) (map[string]interface{}, error)
	UncollectTool(ctx context.Context, userID int, resourceID, resourceType string) (map[string]interface{}, error)
	GetComments(ctx context.Context, resourceID, cursor string, limit int) (map[string]interface{}, error)
	AddComment(ctx context.Context, userID int, resourceID, resourceType, content string) (map[string]interface{}, error)
	DeleteComment(ctx context.Context, userID int, resourceID, commentID string) (map[string]interface{}, error)
	ReplyComment(ctx context.Context, userID int, resourceID, commentID, resourceType, content string) (map[string]interface{}, error)
	DeleteReply(ctx context.Context, userID int, resourceID, commentID string) (map[string]interface{}, error)
	LikeComment(ctx context.Context, userID int, resourceID, commentID string) (map[string]interface{}, error)
	AddView(ctx context.Context, resourceID string) (map[string]interface{}, error)
	UpdateTool(ctx context.Context, userID int, resourceID string, req ToolSubmitRequest) (map[string]interface{}, error)
}

// ToolSubmitRequest 工具提交请求结构体
type ToolSubmitRequest struct {
	Name              string   `form:"name" json:"name" binding:"required"`
	Link              string   `form:"link" json:"link" binding:"required"`
	Description       string   `form:"description" json:"description" binding:"required"`
	DescriptionDetail string   `form:"description_detail" json:"description_detail" binding:"required"`
	Category          string   `form:"catagory" json:"catagory" binding:"required"`
	Tags              []string `form:"tags" json:"tags" binding:"required"`
	Images            []string `form:"images" json:"images"` // 图标/图片URL数组（可选）
	ToolType          string   `form:"type" json:"type"`    // 工具类型：internal/external
}

type toolService struct {
	toolRepo repository.ToolRepository
}

func NewToolService(toolRepo repository.ToolRepository) ToolService {
	return &toolService{toolRepo: toolRepo}
}

func (s *toolService) GetTools(ctx context.Context, category, tags []string, sort, cursor string, pageSize int) (map[string]interface{}, error) {
	tools, err := s.toolRepo.GetTools(ctx, category, tags, sort, cursor, pageSize)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    tools,
	}, nil
}

func (s *toolService) GetTool(ctx context.Context, resourceID, resourceType string, userID int) (map[string]interface{}, error) {
	tool, err := s.toolRepo.GetByID(ctx, resourceID)
	if err != nil {
		return nil, err
	}
	
	// 如果用户已登录，检查点赞和收藏状态
	if userID > 0 {
		isLiked, _ := s.toolRepo.CheckUserLike(ctx, userID, resourceID)
		isCollected, _ := s.toolRepo.CheckUserCollect(ctx, userID, resourceID)
		tool["isliked"] = isLiked
		tool["iscollected"] = isCollected
	}

	return map[string]interface{}{
		"message": "success",
		"data":    tool,
	}, nil
}

func (s *toolService) SearchTools(ctx context.Context, keyword, cursor string, pageSize int, resourceType string) (map[string]interface{}, error) {
	tools, err := s.toolRepo.Search(ctx, keyword, cursor, pageSize)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":     tools,
	}, nil
}

func (s *toolService) SubmitTool(ctx context.Context, userID int, req ToolSubmitRequest) (map[string]interface{}, error) {
	// 处理图片URL（自动本地化外部链接和Base64）
	processedImages := make([]string, 0, len(req.Images))
	for _, imgURL := range req.Images {
		if imgURL != "" {
			// 使用utils.ProcessImageURL自动处理（本地化外部URL和Base64）
			localURL, err := utils.ProcessImageURL(imgURL)
			if err != nil {
				// 如果处理失败，跳过该图片（不强制要求）
				continue
			}
			if localURL != "" {
				processedImages = append(processedImages, localURL)
			}
		}
	}

	// 将结构体转换为 map 传递给 repository
	toolData := map[string]interface{}{
		"name":               req.Name,
		"link":               req.Link,
		"description":        req.Description,
		"description_detail": req.DescriptionDetail,
		"category":           req.Category,
		"tags":               req.Tags,
		"images":             processedImages, // 处理后的本地图片URL数组
		"tool_type":          req.ToolType,    // 工具类型：internal/external
	}

	tool, err := s.toolRepo.Create(ctx, userID, toolData)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "Tool submitted successfully",
		"data":    tool,
	}, nil
}

func (s *toolService) LikeTool(ctx context.Context, userID int, resourceID string) (map[string]interface{}, error) {
	err := s.toolRepo.AddLike(ctx, userID, resourceID)
	if err != nil {
		return nil, err
	}

	likes, _ := s.toolRepo.GetLikes(ctx, resourceID)
	return map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"isliked": true,
			"likes":   likes,
		},
	}, nil
}

func (s *toolService) UnlikeTool(ctx context.Context, userID int, resourceID string) (map[string]interface{}, error) {
	err := s.toolRepo.RemoveLike(ctx, userID, resourceID)
	if err != nil {
		return nil, err
	}
	likes, _ := s.toolRepo.GetLikes(ctx, resourceID)
	return map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"isliked": false,
			"likes":   likes,
		},
	}, nil
}

func (s *toolService) CollectTool(ctx context.Context, userID int, resourceID, resourceType string) (map[string]interface{}, error) {
	err := s.toolRepo.CollectTool(ctx, userID, resourceID)
	if err != nil {
		return nil, err
	}
	
	// 直接从 collections 表统计最新的收藏数，而不是从 tools 表读取
	// 这样可以确保数据的准确性
	collections, err := s.toolRepo.GetCollectionCount(ctx, resourceID)
	if err != nil {
		// 如果获取收藏数失败，返回错误
		return nil, fmt.Errorf("failed to get collection count: %w", err)
	}
	
	return map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"collections": collections,
			"iscollected": true,
		},
	}, nil
}

func (s *toolService) UncollectTool(ctx context.Context, userID int, resourceID, resourceType string) (map[string]interface{}, error) {
	err := s.toolRepo.UncollectTool(ctx, userID, resourceID)
	if err != nil {
		return nil, err
	}
	
	// 直接从 collections 表统计最新的收藏数，而不是从 tools 表读取
	// 这样可以确保数据的准确性
	collections, err := s.toolRepo.GetCollectionCount(ctx, resourceID)
	if err != nil {
		// 如果获取收藏数失败，返回错误
		return nil, fmt.Errorf("failed to get collection count: %w", err)
	}
	
	return map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"collections": collections,
			"iscollected": false,
		},
	}, nil
}

func (s *toolService) GetComments(ctx context.Context, resourceID, cursor string, limit int) (map[string]interface{}, error) {
	cursorInt := 0
	if cursor != "" {
		fmt.Sscanf(cursor, "%d", &cursorInt)
	}
	comments, err := s.toolRepo.GetComments(ctx, resourceID, cursorInt, limit)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"message": "success",
		"data":    comments,
	}, nil
}

func (s *toolService) AddComment(ctx context.Context, userID int, resourceID, resourceType, content string) (map[string]interface{}, error) {
	comment, err := s.toolRepo.AddComment(ctx, userID, resourceID, content)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"message": "success",
		"data":    comment,
	}, nil
}

func (s *toolService) DeleteComment(ctx context.Context, userID int, resourceID, commentID string) (map[string]interface{}, error) {
	err := s.toolRepo.DeleteComment(ctx, userID, resourceID, commentID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"message": "success",
		"data":    map[string]interface{}{},
	}, nil
}

func (s *toolService) ReplyComment(ctx context.Context, userID int, resourceID, commentID, resourceType, content string) (map[string]interface{}, error) {
	reply, err := s.toolRepo.ReplyComment(ctx, userID, resourceID, commentID, content)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"message": "success",
		"data":    reply,
	}, nil
}

func (s *toolService) DeleteReply(ctx context.Context, userID int, resourceID, commentID string) (map[string]interface{}, error) {
	err := s.toolRepo.DeleteReply(ctx, userID, resourceID, commentID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"message": "success",
		"data":    map[string]interface{}{},
	}, nil
}

func (s *toolService) LikeComment(ctx context.Context, userID int, resourceID, commentID string) (map[string]interface{}, error) {
	err := s.toolRepo.LikeComment(ctx, userID, commentID)
	if err != nil {
		return nil, err
	}
	// 获取更新后的点赞数和点赞状态
	likes, _ := s.toolRepo.GetCommentLikes(ctx, commentID)
	isLiked, _ := s.toolRepo.CheckUserCommentLike(ctx, userID, commentID)
	
	return map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"isliked": isLiked,
			"likes":   likes,
		},
	}, nil
}

func (s *toolService) AddView(ctx context.Context, resourceID string) (map[string]interface{}, error) {
	err := s.toolRepo.AddView(ctx, resourceID)
	if err != nil {
		return nil, err
	}
	tool, _ := s.toolRepo.GetByID(ctx, resourceID)
	views := 0
	if tool != nil {
		if v, ok := tool["views"].(int); ok {
			views = v
		}
	}
	return map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"views": views,
		},
	}, nil
}

func (s *toolService) UpdateTool(ctx context.Context, userID int, resourceID string, req ToolSubmitRequest) (map[string]interface{}, error) {
	// 处理图片URL（自动本地化外部链接和Base64）
	processedImages := make([]string, 0, len(req.Images))
	for _, imgURL := range req.Images {
		if imgURL != "" {
			// 使用utils.ProcessImageURL自动处理（本地化外部URL和Base64）
			localURL, err := utils.ProcessImageURL(imgURL)
			if err != nil {
				// 如果处理失败，跳过该图片（不强制要求）
				continue
			}
			if localURL != "" {
				processedImages = append(processedImages, localURL)
			}
		}
	}
	
	// 将结构体转换为 map 传递给 repository
	toolData := map[string]interface{}{
		"name":               req.Name,
		"link":               req.Link,
		"description":        req.Description,
		"description_detail": req.DescriptionDetail,
		"category":           req.Category,
		"tags":               req.Tags,
		"images":             processedImages, // 处理后的本地图片URL数组
		"tool_type":          req.ToolType,    // 工具类型：internal/external
	}
	
	err := s.toolRepo.UpdateTool(ctx, resourceID, userID, toolData)
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"message": "Tool updated successfully",
	}, nil
}
