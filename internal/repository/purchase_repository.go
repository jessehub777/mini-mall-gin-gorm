package repository

import (
	"context"

	"mini-mall-gin-gorm/internal/model"

	"gorm.io/gorm"
)

// PurchaseRepository 封装购买记录数据访问逻辑。
type PurchaseRepository struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) base(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *PurchaseRepository) Create(ctx context.Context, tx *gorm.DB, purchase *model.Purchase) error {
	return r.base(tx).WithContext(ctx).Create(purchase).Error
}

func (r *PurchaseRepository) ListByUserID(ctx context.Context, userID uint) ([]model.Purchase, error) {
	var purchases []model.Purchase
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Product").
		Order("id DESC").
		Find(&purchases).Error
	if err != nil {
		return nil, err
	}
	return purchases, nil
}
