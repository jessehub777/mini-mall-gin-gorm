package model

import "time"

// User 用户表：用于登录认证与个人信息维护。
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Nickname  string    `gorm:"type:varchar(100);not null;default:''" json:"nickname"`
	Email     string    `gorm:"type:varchar(100);not null;default:''" json:"email"`
	Phone     string    `gorm:"type:varchar(30);not null;default:''" json:"phone"`
	Address   string    `gorm:"type:varchar(255);not null;default:''" json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
