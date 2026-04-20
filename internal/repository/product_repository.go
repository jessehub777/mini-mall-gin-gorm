package repository

import (
	"context"

	"mini-mall-gin-gorm/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProductRepository 封装商品数据访问逻辑。
type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) base(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *ProductRepository) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetByIDForUpdate(ctx context.Context, tx *gorm.DB, id uint) (*model.Product, error) {
	var product model.Product
	err := r.base(tx).WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Update(ctx context.Context, tx *gorm.DB, product *model.Product) error {
	return r.base(tx).WithContext(ctx).Save(product).Error
}

func (r *ProductRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, id).Error
}

func (r *ProductRepository) List(ctx context.Context, page, size int) ([]model.Product, int64, error) {
	var (
		products []model.Product
		total    int64
	)

	if err := r.db.WithContext(ctx).Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.db.WithContext(ctx).
		Order("id DESC").
		Offset(offset).
		Limit(size).
		Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
