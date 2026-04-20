package repository

import (
	"context"

	"mini-mall-gin-gorm/internal/model"

	"gorm.io/gorm"
)

// UserRepository 封装用户数据访问逻辑。
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"nickname": user.Nickname,
		"email":    user.Email,
		"phone":    user.Phone,
		"address":  user.Address,
	}).Error
}
