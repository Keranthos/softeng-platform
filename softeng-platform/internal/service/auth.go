package service

import (
	"context"
	"errors"
	"softeng-platform/internal/model"
	"softeng-platform/internal/repository"
	"softeng-platform/internal/utils"
)

type AuthService interface {
	Register(ctx context.Context, req model.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req model.LoginRequest) (*model.User, error)
	ForgotPassword(ctx context.Context, email, newPassword, code string) error
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(ctx context.Context, req model.RegisterRequest) (*model.User, error) {
	// 检查用户名是否已存在
	existingUser, _ := s.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	existingEmail, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	// 验证邮箱验证码和邀请码（这里需要实现具体的验证逻辑）
	if !s.validateEmailCode(req.Email, req.EmailPassword) {
		return nil, errors.New("invalid email verification code")
	}

	if !s.validateCertifyCode(req.CertifyPassword) {
		return nil, errors.New("invalid invitation code")
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Username, // 默认昵称为用户名
		Role:     "user",
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, req model.LoginRequest) (*model.User, error) {
	var user *model.User
	var err error

	// 根据用户名或邮箱查找用户
	if contains(req.UsernameOrEmail, "@") {
		user, err = s.userRepo.GetByEmail(ctx, req.UsernameOrEmail)
	} else {
		user, err = s.userRepo.GetByUsername(ctx, req.UsernameOrEmail)
	}

	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	// 验证密码
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *authService) ForgotPassword(ctx context.Context, email, newPassword, certifyPassword string) error {
	// 验证邀请码（根据API文档，forgot-password使用certify_password作为邀请码）
	if !s.validateCertifyCode(certifyPassword) {
		return errors.New("invalid invitation code")
	}

	// 查找用户
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// 更新密码
	return s.userRepo.UpdatePassword(ctx, user.ID, hashedPassword)
}

func (s *authService) validateEmailCode(email, code string) bool {
	// 实现邮箱验证码验证逻辑
	return true // 临时返回true
}

func (s *authService) validateCertifyCode(code string) bool {
	// 实现邀请码验证逻辑
	return true // 临时返回true
}

func (s *authService) validateResetCode(email, code string) bool {
	// 实现重置密码验证码验证逻辑
	return true // 临时返回true
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
