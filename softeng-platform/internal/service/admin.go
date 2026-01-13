package service

import (
	"context"
	"fmt"
	"softeng-platform/internal/repository"
)

type AdminService interface {
	GetPending(ctx context.Context, itemType string, cursor, limit int, sort string) (map[string]interface{}, error)
	ReviewItem(ctx context.Context, itemID, action, resourceType, rejectReason string) error
}

type adminService struct {
	toolRepo    repository.ToolRepository
	courseRepo  repository.CourseRepository
	projectRepo repository.ProjectRepository
}

func NewAdminService(toolRepo repository.ToolRepository, courseRepo repository.CourseRepository, projectRepo repository.ProjectRepository) AdminService {
	return &adminService{
		toolRepo:    toolRepo,
		courseRepo:  courseRepo,
		projectRepo: projectRepo,
	}
}

func (s *adminService) GetPending(ctx context.Context, itemType string, cursor, limit int, sort string) (map[string]interface{}, error) {
	var data []map[string]interface{}
	var err error

	// 支持前端传递的英文类型名
	switch itemType {
	case "工具", "tools", "tool":
		data, err = s.toolRepo.GetPending(ctx, cursor, limit)
	case "课程", "courses", "course":
		data, err = s.courseRepo.GetPending(ctx, cursor, limit)
	case "项目", "projects", "project":
		data, err = s.projectRepo.GetPending(ctx, cursor, limit)
	case "评论", "comments", "comment":
		// 获取待审核评论
		data = []map[string]interface{}{}
	default:
		// 如果没有指定类型，返回所有类型的待审核项
		toolData, _ := s.toolRepo.GetPending(ctx, cursor, limit)
		courseData, _ := s.courseRepo.GetPending(ctx, cursor, limit)
		projectData, _ := s.projectRepo.GetPending(ctx, cursor, limit)
		data = append(data, toolData...)
		data = append(data, courseData...)
		data = append(data, projectData...)
	}

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":  len(data),
		"cursor": cursor,
		"data":   data,
		"results": data, // 前端可能使用results字段
	}, nil
}

func (s *adminService) ReviewItem(ctx context.Context, itemID, action, resourceType, rejectReason string) error {
	// 根据resourceType判断资源类型
	switch resourceType {
	case "tools", "tool", "工具":
		err := s.toolRepo.UpdateToolStatus(ctx, itemID, action, rejectReason)
		if err != nil {
			return fmt.Errorf("failed to review tool: %w", err)
		}
		return nil
	case "courses", "course", "课程":
		err := s.courseRepo.UpdateCourseStatus(ctx, itemID, action, rejectReason)
		if err != nil {
			return fmt.Errorf("failed to review course: %w", err)
		}
		return nil
	case "projects", "project", "项目":
		err := s.projectRepo.UpdateProjectStatus(ctx, itemID, action, rejectReason)
		if err != nil {
			return fmt.Errorf("failed to review project: %w", err)
		}
		return nil
	default:
		// 如果没有指定类型，尝试工具
		err := s.toolRepo.UpdateToolStatus(ctx, itemID, action, rejectReason)
		if err == nil {
			return nil
		}
		return fmt.Errorf("item not found or cannot be reviewed: %w", err)
	}
}
