package service

import (
	"context"
	"fmt"
	"softeng-platform/internal/model"
	"softeng-platform/internal/repository"
	"softeng-platform/internal/utils"
	"time"
	"sync"
)

type UserService interface {
	GetProfile(ctx context.Context, userID int) (*model.User, error)
	UpdateProfile(ctx context.Context, userID int, req model.UpdateProfileRequest) (*model.User, error)
	GetCollection(ctx context.Context, userID int) (map[string]interface{}, error)
	DeleteCollection(ctx context.Context, userID int, resourceType string, resourceID int) (map[string]interface{}, error)
	GetStatus(ctx context.Context, userID int) (map[string]interface{}, error)
	GetSummit(ctx context.Context, userID int) (map[string]interface{}, error)
	UpdateResourceStatus(ctx context.Context, userID int, resourceType, resourceID, action, state string) (map[string]interface{}, error)
	UpdateEmail(ctx context.Context, userID int, name, password, newEmail, code string) (*model.User, error)
	UpdatePassword(ctx context.Context, userID int, name, email, newPassword, code string) (*model.User, error)
}

type userService struct {
	userRepo    repository.UserRepository
	toolRepo    repository.ToolRepository
	projectRepo repository.ProjectRepository
}

func NewUserService(userRepo repository.UserRepository, toolRepo repository.ToolRepository, projectRepo repository.ProjectRepository) UserService {
	return &userService{userRepo: userRepo, toolRepo: toolRepo, projectRepo: projectRepo}
}

func (s *userService) GetProfile(ctx context.Context, userID int) (*model.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *userService) UpdateProfile(ctx context.Context, userID int, req model.UpdateProfileRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 更新字段（允许设置为空字符串）
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	// Avatar可以为空字符串，表示使用默认头像
	user.Avatar = req.Avatar
	if req.Description != "" {
		user.Description = req.Description
	}
	if req.FacePhoto != "" {
		user.FacePhoto = req.FacePhoto
	}

	err = s.userRepo.Update(ctx, user)
	return user, err
}

func (s *userService) GetCollection(ctx context.Context, userID int) (map[string]interface{}, error) {
	// 获取工具收藏
	toolRows, _ := s.userRepo.GetCollections(ctx, userID, "tool")
	// 获取课程收藏
	courseRows, _ := s.userRepo.GetCollections(ctx, userID, "course")
	// 获取项目收藏
	projectRows, _ := s.userRepo.GetCollections(ctx, userID, "project")
	
	return map[string]interface{}{
		"message":   "success",
		"resources": projectRows,
		"tools":     toolRows,
		"teaches":   courseRows,
	}, nil
}

func (s *userService) DeleteCollection(ctx context.Context, userID int, resourceType string, resourceID int) (map[string]interface{}, error) {
	err := s.userRepo.DeleteCollection(ctx, userID, resourceType, resourceID)
	if err != nil {
		return nil, err
	}
	
	// 重新获取收藏列表
	return s.GetCollection(ctx, userID)
}

func (s *userService) GetStatus(ctx context.Context, userID int) (map[string]interface{}, error) {
	tools, _ := s.userRepo.GetUserPendingItems(ctx, userID, "tool")
	courses, _ := s.userRepo.GetUserPendingItems(ctx, userID, "course")
	projects, _ := s.userRepo.GetUserPendingItems(ctx, userID, "project")
	
	return map[string]interface{}{
		"message":   "success",
		"resources": projects,
		"tools":     tools,
		"teaches":   courses,
	}, nil
}

func (s *userService) GetSummit(ctx context.Context, userID int) (map[string]interface{}, error) {
	// 创建带超时的上下文，每个查询最多10秒
	queryTimeout := 10 * time.Second
	
	// 使用 WaitGroup 并行执行三个查询
	var wg sync.WaitGroup
	var tools, courses, projects []map[string]interface{}
	var toolsErr, coursesErr, projectsErr error
	
	// 查询工具
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[GetSummit] 查询工具时发生panic: %v\n", r)
				toolsErr = fmt.Errorf("panic: %v", r)
				tools = []map[string]interface{}{}
			}
		}()
		queryCtx, cancel := context.WithTimeout(ctx, queryTimeout)
		defer cancel()
		start := time.Now()
		tools, toolsErr = s.userRepo.GetUserSubmissions(queryCtx, userID, "tool")
		elapsed := time.Since(start)
		if toolsErr != nil {
			fmt.Printf("[GetSummit] 查询工具失败 (耗时: %v): %v\n", elapsed, toolsErr)
		} else {
			fmt.Printf("[GetSummit] 查询工具成功 (耗时: %v, 数量: %d)\n", elapsed, len(tools))
		}
	}()
	
	// 查询课程
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[GetSummit] 查询课程时发生panic: %v\n", r)
				coursesErr = fmt.Errorf("panic: %v", r)
				courses = []map[string]interface{}{}
			}
		}()
		queryCtx, cancel := context.WithTimeout(ctx, queryTimeout)
		defer cancel()
		start := time.Now()
		courses, coursesErr = s.userRepo.GetUserSubmissions(queryCtx, userID, "course")
		elapsed := time.Since(start)
		if coursesErr != nil {
			fmt.Printf("[GetSummit] 查询课程失败 (耗时: %v): %v\n", elapsed, coursesErr)
		} else {
			fmt.Printf("[GetSummit] 查询课程成功 (耗时: %v, 数量: %d)\n", elapsed, len(courses))
		}
	}()
	
	// 查询项目
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[GetSummit] 查询项目时发生panic: %v\n", r)
				projectsErr = fmt.Errorf("panic: %v", r)
				projects = []map[string]interface{}{}
			}
		}()
		queryCtx, cancel := context.WithTimeout(ctx, queryTimeout)
		defer cancel()
		start := time.Now()
		projects, projectsErr = s.userRepo.GetUserSubmissions(queryCtx, userID, "project")
		elapsed := time.Since(start)
		if projectsErr != nil {
			fmt.Printf("[GetSummit] 查询项目失败 (耗时: %v): %v\n", elapsed, projectsErr)
		} else {
			fmt.Printf("[GetSummit] 查询项目成功 (耗时: %v, 数量: %d)\n", elapsed, len(projects))
		}
	}()
	
	// 等待所有查询完成
	wg.Wait()
	
	// 如果所有查询都失败，返回错误
	if toolsErr != nil && coursesErr != nil && projectsErr != nil {
		return nil, fmt.Errorf("failed to fetch submissions: tools=%v, courses=%v, projects=%v", toolsErr, coursesErr, projectsErr)
	}
	
	// 如果部分查询失败，记录错误但继续返回成功的结果
	// 将错误结果设置为空数组
	if toolsErr != nil {
		tools = []map[string]interface{}{}
	}
	if coursesErr != nil {
		courses = []map[string]interface{}{}
	}
	if projectsErr != nil {
		projects = []map[string]interface{}{}
	}
	
	return map[string]interface{}{
		"message":   "success",
		"resources": projects,
		"tools":     tools,
		"teaches":   courses,
	}, nil
}

func (s *userService) UpdateResourceStatus(ctx context.Context, userID int, resourceType, resourceID, action, state string) (map[string]interface{}, error) {
	// 验证用户是否有权限操作该资源
	var oldStatus string
	switch resourceType {
	case "tool":
		err := s.userRepo.GetResourceStatus(ctx, resourceType, resourceID, userID)
		if err != nil {
			return nil, err
		}
		oldStatus, _ = s.userRepo.GetResourceOldStatus(ctx, resourceType, resourceID)
		err = s.toolRepo.UpdateToolStatus(ctx, resourceID, action, "")
		if err != nil {
			return nil, err
		}
		s.userRepo.LogStatusChange(ctx, resourceType, resourceID, oldStatus, action, userID)
	case "project":
		// 项目支持状态更新
		err := s.userRepo.GetResourceStatus(ctx, resourceType, resourceID, userID)
		if err != nil {
			return nil, err
		}
		oldStatus, _ = s.userRepo.GetResourceOldStatus(ctx, resourceType, resourceID)
		err = s.projectRepo.UpdateProjectStatus(ctx, resourceID, action, "")
		if err != nil {
			return nil, err
		}
		s.userRepo.LogStatusChange(ctx, resourceType, resourceID, oldStatus, action, userID)
	case "course":
		// 课程暂时不支持状态更新（没有status字段）
		return nil, fmt.Errorf("status update not supported for %s", resourceType)
	default:
		return nil, fmt.Errorf("invalid resource type")
	}
	
	return map[string]interface{}{
		"message": "success",
		"manipulate": map[string]interface{}{
			"resourceId":   resourceID,
			"resourceType": resourceType,
			"newstatus":    action,
			"oldestatus":   oldStatus,
			"operateTime":  time.Now().Format("2006-01-02 15:04:05"),
			"operator":     "user",
		},
	}, nil
}

func (s *userService) UpdateEmail(ctx context.Context, userID int, name, password, newEmail, code string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	// 验证用户名和密码
	if user.Username != name {
		return nil, fmt.Errorf("invalid username")
	}
	
	// 验证密码（这里应该使用密码验证函数，暂时简化）
	// if !utils.CheckPasswordHash(password, user.Password) {
	//     return nil, fmt.Errorf("invalid password")
	// }
	
	// TODO: 验证邮箱验证码 code
	
	// 检查新邮箱是否已存在
	existingUser, _ := s.userRepo.GetByEmail(ctx, newEmail)
	if existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}
	
	// 更新邮箱
	user.Email = newEmail
	err = s.userRepo.UpdateEmail(ctx, userID, newEmail)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (s *userService) UpdatePassword(ctx context.Context, userID int, name, email, newPassword, code string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	// 验证用户名和邮箱
	if user.Username != name || user.Email != email {
		return nil, fmt.Errorf("invalid username or email")
	}
	
	// TODO: 验证密码重置验证码 code
	
	// 加密新密码并更新
	hashedPassword, err := s.hashPassword(newPassword)
	if err != nil {
		return nil, err
	}
	
	err = s.userRepo.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (s *userService) hashPassword(password string) (string, error) {
	// 使用bcrypt加密密码
	return utils.HashPassword(password)
}
