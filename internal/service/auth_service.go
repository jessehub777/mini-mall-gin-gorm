package service

import (
	"context"
	"errors"

	"mini-mall-gin-gorm/internal/config"
	"mini-mall-gin-gorm/internal/model"
	"mini-mall-gin-gorm/internal/repository"
	"mini-mall-gin-gorm/pkg/jwtutil"
	"mini-mall-gin-gorm/pkg/password"

	"gorm.io/gorm"
)

var (
	// ErrUserExists 表示注册时用户名已存在。
	ErrUserExists = errors.New("用户名已存在")
	// ErrInvalidCredentials 表示用户名或密码错误。
	ErrInvalidCredentials = errors.New("用户名或密码错误")
)

// RegisterRequest 定义注册请求字段。
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Nickname string `json:"nickname" binding:"omitempty,max=100"`
}

// LoginRequest 定义登录请求字段。
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthService 处理认证相关业务。
type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*model.User, error) {
	// 先查重，避免重复用户名。
	_, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, ErrUserExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashed, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Password: hashed,
		Nickname: req.Nickname,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (string, *model.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}

	if !password.Verify(req.Password, user.Password) {
		return "", nil, ErrInvalidCredentials
	}

	token, err := jwtutil.GenerateToken(s.cfg.JWT.Secret, s.cfg.JWT.ExpireHours, user.ID, user.Username)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
