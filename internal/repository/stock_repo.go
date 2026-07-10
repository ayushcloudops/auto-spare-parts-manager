package repository

import (
	"context"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"

	"gorm.io/gorm"
)

// stockRepo is the GORM implementation of domain.StockRepository.
type stockRepo struct {
	Base[domain.StockMovement]
}

// NewStockRepo constructs a StockRepository.
func NewStockRepo(db *gorm.DB) domain.StockRepository {
	return &stockRepo{Base: NewBase[domain.StockMovement](db)}
}

// Record appends a stock movement to the ledger.
func (r *stockRepo) Record(ctx context.Context, m *domain.StockMovement) error {
	return r.Create(ctx, m)
}

// ListByProduct returns a product's movements, newest first.
func (r *stockRepo) ListByProduct(ctx context.Context, productID uint) ([]domain.StockMovement, error) {
	var out []domain.StockMovement
	err := r.Conn(ctx).
		Where("product_id = ?", productID).
		Order("occurred_at DESC").
		Find(&out).Error
	if err != nil {
		return nil, apperr.Internal(err, "could not load stock history")
	}
	return out, nil
}
