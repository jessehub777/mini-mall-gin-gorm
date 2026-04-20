package service

import (
	"context"
	"errors"

	"mini-mall-gin-gorm/internal/model"
	"mini-mall-gin-gorm/internal/repository"

	"gorm.io/gorm"
)

var (
	// ErrProductNotFound 表示商品不存在。
	ErrProductNotFound = errors.New("商品不存在")
	// ErrInvalidPrice 表示价格非法。
	ErrInvalidPrice = errors.New("价格不能小于 0")
	// ErrInvalidStock 表示库存非法。
	ErrInvalidStock = errors.New("库存不能小于 0")
)

// CreateProductRequest 定义创建商品请求字段。
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,max=120"`
	Description string  `json:"description" binding:"omitempty,max=500"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int     `json:"stock" binding:"required"`
}

// UpdateProductRequest 定义更新商品请求字段。
type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"required,max=120"`
	Description string  `json:"description" binding:"omitempty,max=500"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int     `json:"stock" binding:"required"`
}

// ProductService 处理商品业务。
type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) Create(ctx context.Context, req CreateProductRequest) (*model.Product, error) {
	if req.Price < 0 {
		return nil, ErrInvalidPrice
	}
	if req.Stock < 0 {
		return nil, ErrInvalidStock
	}

	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}
	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) Update(ctx context.Context, id uint, req UpdateProductRequest) (*model.Product, error) {
	if req.Price < 0 {
		return nil, ErrInvalidPrice
	}
	if req.Stock < 0 {
		return nil, ErrInvalidStock
	}

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock

	if err := s.productRepo.Update(ctx, nil, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}
	return s.productRepo.Delete(ctx, product.ID)
}

func (s *ProductService) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

func (s *ProductService) List(ctx context.Context, page, size int) ([]model.Product, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return s.productRepo.List(ctx, page, size)
}
