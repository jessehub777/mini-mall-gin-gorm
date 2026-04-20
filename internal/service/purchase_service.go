package service

import (
	"context"
	"errors"

	"mini-mall-gin-gorm/internal/model"
	"mini-mall-gin-gorm/internal/repository"

	"gorm.io/gorm"
)

var (
	// ErrInvalidQuantity 表示购买数量不合法。
	ErrInvalidQuantity = errors.New("购买数量必须大于 0")
	// ErrInsufficientStock 表示库存不足。
	ErrInsufficientStock = errors.New("库存不足")
)

// CreatePurchaseRequest 定义购买请求字段。
type CreatePurchaseRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required"`
}

// PurchaseService 处理购买交易逻辑。
type PurchaseService struct {
	db           *gorm.DB
	productRepo  *repository.ProductRepository
	purchaseRepo *repository.PurchaseRepository
}

func NewPurchaseService(db *gorm.DB, productRepo *repository.ProductRepository, purchaseRepo *repository.PurchaseRepository) *PurchaseService {
	return &PurchaseService{db: db, productRepo: productRepo, purchaseRepo: purchaseRepo}
}

func (s *PurchaseService) Create(ctx context.Context, userID uint, req CreatePurchaseRequest) (*model.Purchase, error) {
	if req.Quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	var result *model.Purchase
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 通过行锁保证并发购买时库存扣减正确。
		product, err := s.productRepo.GetByIDForUpdate(ctx, tx, req.ProductID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrProductNotFound
			}
			return err
		}

		if product.Stock < req.Quantity {
			return ErrInsufficientStock
		}

		product.Stock -= req.Quantity
		if err := s.productRepo.Update(ctx, tx, product); err != nil {
			return err
		}

		purchase := &model.Purchase{
			UserID:      userID,
			ProductID:   product.ID,
			Quantity:    req.Quantity,
			UnitPrice:   product.Price,
			TotalAmount: product.Price * float64(req.Quantity),
		}
		if err := s.purchaseRepo.Create(ctx, tx, purchase); err != nil {
			return err
		}

		result = purchase
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PurchaseService) ListByUserID(ctx context.Context, userID uint) ([]model.Purchase, error) {
	return s.purchaseRepo.ListByUserID(ctx, userID)
}
