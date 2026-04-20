package service

import (
	"context"

	"mini-mall-gin-gorm/internal/model"
	"mini-mall-gin-gorm/internal/repository"
)

// UpdateProfileRequest 定义更新个人信息请求字段。
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,max=100"`
	Email    string `json:"email" binding:"omitempty,max=100"`
	Phone    string `json:"phone" binding:"omitempty,max=30"`
	Address  string `json:"address" binding:"omitempty,max=255"`
}

// UserService 处理用户资料业务。
type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(ctx context.Context, userID uint) (*model.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uint, req UpdateProfileRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Phone = req.Phone
	user.Address = req.Address

	if err := s.userRepo.UpdateProfile(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
