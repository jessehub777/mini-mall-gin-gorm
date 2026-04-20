package model

import "time"

// Purchase 用户购买记录表：记录谁在何时购买了什么、买了多少。
type Purchase struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	ProductID   uint      `gorm:"not null;index" json:"product_id"`
	Quantity    int       `gorm:"not null" json:"quantity"`
	UnitPrice   float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	TotalAmount float64   `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`

	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (Purchase) TableName() string {
	return "purchases"
}
