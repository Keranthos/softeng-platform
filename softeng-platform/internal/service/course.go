package service

import (
	"context"
	"softeng-platform/internal/repository"
)

type CourseService interface {
	GetCourses(ctx context.Context, semester string, category []string, sort string, limit, cursor int, resourceType string) (map[string]interface{}, error)
	GetCourse(ctx context.Context, courseID, resourceType string) (map[string]interface{}, error)
	SearchCourses(ctx context.Context, keyword string, category []string, limit, cursor int, resourceType string) (map[string]interface{}, error)
	UploadResource(ctx context.Context, userID int, courseID, resourceType string, req CourseUploadRequest) (map[string]interface{}, error)
	DownloadTextbook(ctx context.Context, courseID, textbookID string) (map[string]interface{}, error)
	AddComment(ctx context.Context, userID int, courseID, content string) (map[string]interface{}, error)
	DeleteComment(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
	ReplyComment(ctx context.Context, userID int, courseID, commentID, content string) (map[string]interface{}, error)
	DeleteReply(ctx context.Context, userID int, courseID, commentID string) (map[string]interface{}, error)
	AddView(ctx context.Context, courseID string) (map[string]interface{}, error)
	CollectCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
	UncollectCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
	LikeCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
	UnlikeCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
}

// CourseUploadRequest 课程资源上传请求
type CourseUploadRequest struct {
	File        string   `form:"file" json:"file"`
	Resource    string   `form:"resource" json:"resource"`
	Description string   `form:"description" json:"description" binding:"required"`
	Tags        []string `form:"tags" json:"tags"`
}

type courseService struct {
	courseRepo repository.CourseRepository
}

func NewCourseService(courseRepo repository.CourseRepository) CourseService {
	return &courseService{courseRepo: courseRepo}
}

func (s *courseService) GetCourses(ctx context.Context, semester string, category []string, sort string, limit, cursor int, resourceType string) (map[string]interface{}, error) {
	courses, err := s.courseRepo.GetCourses(ctx, semester, category, sort, limit, cursor)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message":     "success",
		"courses_agg": courses,
	}, nil
}

func (s *courseService) GetCourse(ctx context.Context, courseID, resourceType string) (map[string]interface{}, error) {
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"courses": []map[string]interface{}{course},
	}, nil
}

func (s *courseService) SearchCourses(ctx context.Context, keyword string, category []string, limit, cursor int, resourceType string) (map[string]interface{}, error) {
	courses, err := s.courseRepo.Search(ctx, keyword, category, limit, cursor)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message":     "success",
		"courses_agg": courses,
	}, nil
}

func (s *courseService) UploadResource(ctx context.Context, userID int, courseID, resourceType string, req CourseUploadRequest) (map[string]interface{}, error) {
	// 将结构体转换为 map 传递给 repository
	resourceData := map[string]interface{}{
		"file":        req.File,
		"resource":    req.Resource,
		"description": req.Description,
		"tags":        req.Tags,
	}

	resource, err := s.courseRepo.UploadResource(ctx, userID, courseID, resourceData)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "Resource uploaded successfully",
		"data":    resource,
	}, nil
}

// 其他方法保持不变...
func (s *courseService) DownloadTextbook(ctx context.Context, courseID, textbookID string) (map[string]interface{}, error) {
	content, err := s.courseRepo.DownloadTextbook(ctx, courseID, textbookID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"content": content,
	}, nil
}

func (s *courseService) AddComment(ctx context.Context, userID int, courseID, content string) (map[string]interface{}, error) {
	comment, err := s.courseRepo.AddComment(ctx, userID, courseID, content)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    comment,
	}, nil
}

func (s *courseService) DeleteComment(ctx context.Context, userID int, courseID string) (map[string]interface{}, error) {
	comment, err := s.courseRepo.DeleteComment(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    comment,
	}, nil
}

func (s *courseService) ReplyComment(ctx context.Context, userID int, courseID, commentID, content string) (map[string]interface{}, error) {
	reply, err := s.courseRepo.ReplyComment(ctx, userID, courseID, commentID, content)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    reply,
	}, nil
}

func (s *courseService) DeleteReply(ctx context.Context, userID int, courseID, commentID string) (map[string]interface{}, error) {
	reply, err := s.courseRepo.DeleteReply(ctx, userID, courseID, commentID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    reply,
	}, nil
}

func (s *courseService) AddView(ctx context.Context, courseID string) (map[string]interface{}, error) {
	views, err := s.courseRepo.AddView(ctx, courseID)
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

func (s *courseService) CollectCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error) {
	result, err := s.courseRepo.CollectCourse(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    result,
	}, nil
}

func (s *courseService) UncollectCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error) {
	result, err := s.courseRepo.UncollectCourse(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    result,
	}, nil
}

func (s *courseService) LikeCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error) {
	result, err := s.courseRepo.LikeCourse(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    result,
	}, nil
}

func (s *courseService) UnlikeCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error) {
	result, err := s.courseRepo.UnlikeCourse(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    result,
	}, nil
}
