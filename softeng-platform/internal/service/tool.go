package service

import (
	"context"
	"softeng-platform/internal/repository"
)

type ToolService interface {
	GetTools(ctx context.Context, category, tags []string, sort, cursor string, pageSize int) (map[string]interface{}, error)
	GetTool(ctx context.Context, resourceID, resourceType string) (map[string]interface{}, error)
	SearchTools(ctx context.Context, keyword, cursor string, pageSize int, resourceType string) (map[string]interface{}, error)
	SubmitTool(ctx context.Context, userID int, req ToolSubmitRequest) (map[string]interface{}, error)
	LikeTool(ctx context.Context, userID int, resourceID string) (map[string]interface{}, error)
	UnlikeTool(ctx context.Context, userID int, resourceID string) (map[string]interface{}, error)
	CollectTool(ctx context.Context, userID int, resourceID, resourceType string) (map[string]interface{}, error)
	UncollectTool(ctx context.Context, userID int, resourceID, resourceType string) (map[string]interface{}, error)
	AddComment(ctx context.Context, userID int, resourceID, resourceType, content string) (map[string]interface{}, error)
	DeleteComment(ctx context.Context, userID int, resourceID string) (map[string]interface{}, error)
	ReplyComment(ctx context.Context, userID int, resourceID, commentID, resourceType, content string) (map[string]interface{}, error)
	DeleteReply(ctx context.Context, userID int, resourceID, commentID string) (map[string]interface{}, error)
	AddView(ctx context.Context, resourceID string) (map[string]interface{}, error)
}

// ToolSubmitRequest 工具提交请求结构体
type ToolSubmitRequest struct {
	Name              string   `form:"name" json:"name" binding:"required"`
	Link              string   `form:"link" json:"link" binding:"required"`
	Description       string   `form:"description" json:"description" binding:"required"`
	DescriptionDetail string   `form:"description_detail" json:"description_detail" binding:"required"`
	Category          string   `form:"catagory" json:"catagory" binding:"required"`
	Tags              []string `form:"tags" json:"tags" binding:"required"`
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

func (s *toolService) GetTool(ctx context.Context, resourceID, resourceType string) (map[string]interface{}, error) {
	tool, err := s.toolRepo.GetByID(ctx, resourceID)
	if err != nil {
		return nil, err
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
		"messages": "success",
		"data":     tools,
	}, nil
}

func (s *toolService) SubmitTool(ctx context.Context, userID int, req ToolSubmitRequest) (map[string]interface{}, error) {
	// 将结构体转换为 map 传递给 repository
	toolData := map[string]interface{}{
		"name":               req.Name,
		"link":               req.Link,
		"description":        req.Description,
		"description_detail": req.DescriptionDetail,
		"category":           req.Category,
		"tags":               req.Tags,
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

// 其他方法保持不变...
func (s *toolService) UnlikeTool(ctx context.Context, userID int, resourceID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"isliked": false,
			"likes":   0,
		},
	}, nil
}

func (s *toolService) CollectTool(ctx context.Context, userID int, resourceID, resourceType string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"messages": "success",
		"data": map[string]interface{}{
			"collections": 0,
			"iscollected": true,
		},
	}, nil
}

func (s *toolService) UncollectTool(ctx context.Context, userID int, resourceID, resourceType string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"messages": "success",
		"data": map[string]interface{}{
			"collections": 0,
			"iscollected": false,
		},
	}, nil
}

func (s *toolService) AddComment(ctx context.Context, userID int, resourceID, resourceType, content string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"message": "success",
		"com":     map[string]interface{}{},
	}, nil
}

func (s *toolService) DeleteComment(ctx context.Context, userID int, resourceID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"message": "success",
		"data":    map[string]interface{}{},
	}, nil
}

func (s *toolService) ReplyComment(ctx context.Context, userID int, resourceID, commentID, resourceType, content string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"message": "success",
		"data":    map[string]interface{}{},
	}, nil
}

func (s *toolService) DeleteReply(ctx context.Context, userID int, resourceID, commentID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"message": "success",
		"data":    map[string]interface{}{},
	}, nil
}

func (s *toolService) AddView(ctx context.Context, resourceID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"views": 0,
		},
	}, nil
}
