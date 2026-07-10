package repository

import (
	"context"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"
	"autoshop/internal/pkg/money"

	"gorm.io/gorm"
)

type customerRepo struct {
	Base[domain.Customer]
}

// NewCustomerRepo constructs a CustomerRepository.
func NewCustomerRepo(db *gorm.DB) domain.CustomerRepository {
	return &customerRepo{Base: NewBase[domain.Customer](db)}
}

func (r *customerRepo) List(ctx context.Context, search string) ([]domain.Customer, error) {
	q := r.Conn(ctx).Model(&domain.Customer{}).Order("name ASC")
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("name LIKE ? OR phone LIKE ?", like, like)
	}
	var out []domain.Customer
	if err := q.Find(&out).Error; err != nil {
		return nil, apperr.Internal(err, "could not list customers")
	}
	return out, nil
}

// AdjustOutstanding atomically adds delta to the customer's balance using a SQL
// expression, avoiding a read-modify-write race.
func (r *customerRepo) AdjustOutstanding(ctx context.Context, id uint, delta money.Money) error {
	res := r.Conn(ctx).Model(&domain.Customer{}).
		Where("id = ?", id).
		UpdateColumn("outstanding", gorm.Expr("outstanding + ?", delta.Paise()))
	if res.Error != nil {
		return apperr.Internal(res.Error, "could not update outstanding")
	}
	if res.RowsAffected == 0 {
		return apperr.NotFound("customer #%d not found", id)
	}
	return nil
}

func (r *customerRepo) SumOutstanding(ctx context.Context) (money.Money, error) {
	var total int64
	err := r.Conn(ctx).Model(&domain.Customer{}).
		Select("COALESCE(SUM(outstanding), 0)").Scan(&total).Error
	if err != nil {
		return 0, apperr.Internal(err, "could not sum outstanding")
	}
	return money.FromPaise(total), nil
}
