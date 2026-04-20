package model

import "time"

// Product 商品表：用于展示商品与管理库存。
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(120);not null" json:"name"`
	Description string    `gorm:"type:varchar(500);not null;default:''" json:"description"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int       `gorm:"not null;default:0" json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}
