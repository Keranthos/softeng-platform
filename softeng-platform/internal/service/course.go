package service

import (
	"context"
	"softeng-platform/internal/repository"
)

type CourseService interface {
	GetCourses(ctx context.Context, semester string, category []string, sort string, limit, cursor int, resourceType string) (map[string]interface{}, error)
	GetCourse(ctx context.Context, courseID, resourceType string, userID int) (map[string]interface{}, error)
	SearchCourses(ctx context.Context, keyword string, category []string, limit, cursor int, resourceType string) (map[string]interface{}, error)
	SubmitCourse(ctx context.Context, userID int, req CourseSubmitRequest) (map[string]interface{}, error)
	UploadResource(ctx context.Context, userID int, courseID, resourceType string, req CourseUploadRequest) (map[string]interface{}, error)
	DownloadTextbook(ctx context.Context, courseID, textbookID string) (map[string]interface{}, error)
	GetComments(ctx context.Context, courseID string, cursor, limit int) (map[string]interface{}, error)
	AddComment(ctx context.Context, userID int, courseID, content string) (map[string]interface{}, error)
	DeleteComment(ctx context.Context, userID int, courseID, commentID string) (map[string]interface{}, error)
	ReplyComment(ctx context.Context, userID int, courseID, commentID, content string) (map[string]interface{}, error)
	DeleteReply(ctx context.Context, userID int, courseID, commentID string) (map[string]interface{}, error)
	LikeComment(ctx context.Context, userID int, courseID, commentID string) (map[string]interface{}, error)
	GetResources(ctx context.Context, courseID string) (map[string]interface{}, error)
	AddView(ctx context.Context, courseID string) (map[string]interface{}, error)
	CollectCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
	UncollectCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
	LikeCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
	UnlikeCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
}

// CourseSubmitRequest 课程提交请求
type CourseSubmitRequest struct {
	Name      string   `form:"name" json:"name" binding:"required"`
	Semester  string   `form:"semester" json:"semester"`
	Credit    int      `form:"credit" json:"credit"`
	Teacher   []string `form:"teacher" json:"teacher"`
	Category  []string `form:"category" json:"category"`
	Cover     string   `form:"cover" json:"cover"`
	Resources []string `form:"resources" json:"resources"`
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

func (s *courseService) GetCourse(ctx context.Context, courseID, resourceType string, userID int) (map[string]interface{}, error) {
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	
	// 如果用户已登录，检查点赞和收藏状态
	if userID > 0 {
		isLiked, _ := s.courseRepo.CheckUserLike(ctx, userID, courseID)
		isCollected, _ := s.courseRepo.CheckUserCollect(ctx, userID, courseID)
		course["isliked"] = isLiked
		course["iscollected"] = isCollected
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

func (s *courseService) DeleteComment(ctx context.Context, userID int, courseID, commentID string) (map[string]interface{}, error) {
	comment, err := s.courseRepo.DeleteComment(ctx, userID, courseID, commentID)
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

// SubmitCourse 提交课程
func (s *courseService) SubmitCourse(ctx context.Context, userID int, req CourseSubmitRequest) (map[string]interface{}, error) {
	courseData := map[string]interface{}{
		"name":      req.Name,
		"semester":  req.Semester,
		"credit":    req.Credit,
		"teacher":   req.Teacher,
		"category":  req.Category,
		"cover":     req.Cover,
		"resources": req.Resources,
	}

	course, err := s.courseRepo.Create(ctx, userID, courseData)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "Course submitted successfully",
		"data":    course,
	}, nil
}

// GetComments 获取课程评论列表
func (s *courseService) GetComments(ctx context.Context, courseID string, cursor, limit int) (map[string]interface{}, error) {
	comments, err := s.courseRepo.GetComments(ctx, courseID, cursor, limit)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    comments,
	}, nil
}

// LikeComment 点赞/取消点赞评论
func (s *courseService) LikeComment(ctx context.Context, userID int, courseID, commentID string) (map[string]interface{}, error) {
	result, err := s.courseRepo.LikeComment(ctx, userID, courseID, commentID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    result,
	}, nil
}

// GetResources 获取课程资源
func (s *courseService) GetResources(ctx context.Context, courseID string) (map[string]interface{}, error) {
	resources, err := s.courseRepo.GetResources(ctx, courseID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "success",
		"data":    resources,
	}, nil
}
