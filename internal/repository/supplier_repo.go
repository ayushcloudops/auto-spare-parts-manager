package repository

import (
	"context"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"

	"gorm.io/gorm"
)

type supplierRepo struct {
	Base[domain.Supplier]
}

// NewSupplierRepo constructs a SupplierRepository.
func NewSupplierRepo(db *gorm.DB) domain.SupplierRepository {
	return &supplierRepo{Base: NewBase[domain.Supplier](db)}
}

func (r *supplierRepo) List(ctx context.Context, search string) ([]domain.Supplier, error) {
	q := r.Conn(ctx).Model(&domain.Supplier{}).Order("name ASC")
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("name LIKE ? OR phone LIKE ?", like, like)
	}
	var out []domain.Supplier
	if err := q.Find(&out).Error; err != nil {
		return nil, apperr.Internal(err, "could not list suppliers")
	}
	return out, nil
}
