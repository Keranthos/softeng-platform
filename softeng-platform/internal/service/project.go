package service

import (
	"context"
	"softeng-platform/internal/repository"
)

type ProjectService interface {
	GetProjects(ctx context.Context, category string, techStack []string, sort string, limit int, cursor, resourceType string) (map[string]interface{}, error)
	GetProject(ctx context.Context, projectID string) (map[string]interface{}, error)
	SearchProjects(ctx context.Context, keyword string, category []string, cursor string, limit int) (map[string]interface{}, error)
	UploadProject(ctx context.Context, userID int, req ProjectUploadRequest) (map[string]interface{}, error)
	UpdateProject(ctx context.Context, userID int, projectID string, req ProjectUploadRequest) (map[string]interface{}, error)
	LikeProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error)
	UnlikeProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error)
	AddComment(ctx context.Context, userID int, projectID, content string) (map[string]interface{}, error)
	DeleteComment(ctx context.Context, userID int, projectID string) (map[string]interface{}, error)
	ReplyComment(ctx context.Context, userID int, projectID, commentID, content string) (map[string]interface{}, error)
	DeleteReply(ctx context.Context, userID int, projectID, commentID string) (map[string]interface{}, error)
	AddView(ctx context.Context, projectID string) (map[string]interface{}, error)
	CollectProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error)
	UncollectProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error)
}

// ProjectUploadRequest 项目上传请求
type ProjectUploadRequest struct {
	Name        string   `form:"name" json:"name" binding:"required"`
	Description string   `form:"description" json:"description" binding:"required"`
	Detail      string   `form:"detail" json:"detail" binding:"required"`
	Github      string   `form:"github" json:"github"`
	TechStack   []string `form:"techStack" json:"techStack" binding:"required"`
	Category    string   `form:"catagory" json:"catagory" binding:"required"`
	Images      []string `form:"images" json:"images"`
}

type projectService struct {
	projectRepo repository.ProjectRepository
}

func NewProjectService(projectRepo repository.ProjectRepository) ProjectService {
	return &projectService{projectRepo: projectRepo}
}

func (s *projectService) GetProjects(ctx context.Context, category string, techStack []string, sort string, limit int, cursor, resourceType string) (map[string]interface{}, error) {
	projects, err := s.projectRepo.GetProjects(ctx, category, techStack, sort, limit, cursor)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    projects,
	}, nil
}

func (s *projectService) GetProject(ctx context.Context, projectID string) (map[string]interface{}, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    project,
	}, nil
}

func (s *projectService) SearchProjects(ctx context.Context, keyword string, category []string, cursor string, limit int) (map[string]interface{}, error) {
	projects, err := s.projectRepo.Search(ctx, keyword, category, cursor, limit)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    projects,
	}, nil
}

func (s *projectService) UploadProject(ctx context.Context, userID int, req ProjectUploadRequest) (map[string]interface{}, error) {
	// 将结构体转换为 map 传递给 repository
	projectData := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"detail":      req.Detail,
		"github":      req.Github,
		"techStack":   req.TechStack,
		"category":    req.Category,
		"images":      req.Images,
	}

	project, err := s.projectRepo.Create(ctx, userID, projectData)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "Project uploaded successfully",
		"data":    project,
	}, nil
}

func (s *projectService) UpdateProject(ctx context.Context, userID int, projectID string, req ProjectUploadRequest) (map[string]interface{}, error) {
	// 将结构体转换为 map 传递给 repository
	projectData := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"detail":      req.Detail,
		"github":      req.Github,
		"techStack":   req.TechStack,
		"category":    req.Category,
		"images":      req.Images,
	}

	project, err := s.projectRepo.Update(ctx, userID, projectID, projectData)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "Project updated successfully",
		"data":    project,
	}, nil
}

func (s *projectService) LikeProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error) {
	result, err := s.projectRepo.LikeProject(ctx, userID, projectID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    result,
	}, nil
}

func (s *projectService) UnlikeProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error) {
	result, err := s.projectRepo.UnlikeProject(ctx, userID, projectID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    result,
	}, nil
}

func (s *projectService) AddComment(ctx context.Context, userID int, projectID, content string) (map[string]interface{}, error) {
	comment, err := s.projectRepo.AddComment(ctx, userID, projectID, content)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    comment,
	}, nil
}

func (s *projectService) DeleteComment(ctx context.Context, userID int, projectID string) (map[string]interface{}, error) {
	comment, err := s.projectRepo.DeleteComment(ctx, userID, projectID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    comment,
	}, nil
}

func (s *projectService) ReplyComment(ctx context.Context, userID int, projectID, commentID, content string) (map[string]interface{}, error) {
	reply, err := s.projectRepo.ReplyComment(ctx, userID, projectID, commentID, content)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    reply,
	}, nil
}

func (s *projectService) DeleteReply(ctx context.Context, userID int, projectID, commentID string) (map[string]interface{}, error) {
	reply, err := s.projectRepo.DeleteReply(ctx, userID, projectID, commentID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    reply,
	}, nil
}

func (s *projectService) AddView(ctx context.Context, projectID string) (map[string]interface{}, error) {
	views, err := s.projectRepo.AddView(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"views": views,
		},
	}, nil
}

func (s *projectService) CollectProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error) {
	result, err := s.projectRepo.CollectProject(ctx, userID, projectID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    result,
	}, nil
}

func (s *projectService) UncollectProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error) {
	result, err := s.projectRepo.UncollectProject(ctx, userID, projectID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    result,
	}, nil
}
