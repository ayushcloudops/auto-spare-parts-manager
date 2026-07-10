package repository

import (
	"context"
	"errors"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"

	"gorm.io/gorm"
)

type purchaseRepo struct {
	Base[domain.Purchase]
}

// NewPurchaseRepo constructs a PurchaseRepository.
func NewPurchaseRepo(db *gorm.DB) domain.PurchaseRepository {
	return &purchaseRepo{Base: NewBase[domain.Purchase](db)}
}

func (r *purchaseRepo) FindByID(ctx context.Context, id uint) (*domain.Purchase, error) {
	var p domain.Purchase
	err := r.Conn(ctx).Preload("Items").Preload("Supplier").First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.NotFound("purchase #%d not found", id)
	}
	if err != nil {
		return nil, apperr.Internal(err, "could not load purchase")
	}
	return &p, nil
}

func (r *purchaseRepo) List(ctx context.Context, f domain.PurchaseFilter) ([]domain.Purchase, error) {
	q := r.Conn(ctx).Model(&domain.Purchase{}).Preload("Supplier").Order("date DESC, id DESC")
	if f.SupplierID != nil {
		q = q.Where("supplier_id = ?", *f.SupplierID)
	}
	if f.Limit > 0 {
		q = q.Limit(f.Limit).Offset(f.Offset)
	}
	var out []domain.Purchase
	if err := q.Find(&out).Error; err != nil {
		return nil, apperr.Internal(err, "could not list purchases")
	}
	return out, nil
}
