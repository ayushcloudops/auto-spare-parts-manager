package repository

import (
	"context"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"

	"gorm.io/gorm"
)

// productRepo is the GORM implementation of domain.ProductRepository. It embeds
// Base[Product] for the standard CRUD and adds product-specific queries.
type productRepo struct {
	Base[domain.Product]
}

// NewProductRepo constructs a ProductRepository.
func NewProductRepo(db *gorm.DB) domain.ProductRepository {
	return &productRepo{Base: NewBase[domain.Product](db)}
}

// applyFilter turns a ProductFilter into WHERE clauses. Zero-value fields are
// ignored so one method serves list/search/low-stock.
func applyProductFilter(q *gorm.DB, f domain.ProductFilter) *gorm.DB {
	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where(
			"name LIKE ? OR part_number LIKE ? OR brand LIKE ?",
			like, like, like,
		)
	}
	if f.Category != "" {
		q = q.Where("category = ?", f.Category)
	}
	if f.VehicleBrand != "" {
		q = q.Where("vehicle_brand = ?", f.VehicleBrand)
	}
	if f.LowStockOnly {
		q = q.Where("current_stock <= minimum_stock")
	}
	return q
}

func (r *productRepo) List(ctx context.Context, f domain.ProductFilter) ([]domain.Product, error) {
	q := applyProductFilter(r.Conn(ctx).Model(&domain.Product{}), f).Order("name ASC")
	if f.Limit > 0 {
		q = q.Limit(f.Limit).Offset(f.Offset)
	}
	var out []domain.Product
	if err := q.Find(&out).Error; err != nil {
		return nil, apperr.Internal(err, "could not list products")
	}
	return out, nil
}

func (r *productRepo) Count(ctx context.Context, f domain.ProductFilter) (int64, error) {
	var n int64
	q := applyProductFilter(r.Conn(ctx).Model(&domain.Product{}), f)
	if err := q.Count(&n).Error; err != nil {
		return 0, apperr.Internal(err, "could not count products")
	}
	return n, nil
}

func (r *productRepo) CountAll(ctx context.Context) (int64, error) {
	var n int64
	if err := r.Conn(ctx).Model(&domain.Product{}).Count(&n).Error; err != nil {
		return 0, apperr.Internal(err, "could not count products")
	}
	return n, nil
}

func (r *productRepo) CountLowStock(ctx context.Context) (int64, error) {
	var n int64
	err := r.Conn(ctx).Model(&domain.Product{}).
		Where("current_stock <= minimum_stock").Count(&n).Error
	if err != nil {
		return 0, apperr.Internal(err, "could not count low stock")
	}
	return n, nil
}

func (r *productRepo) Categories(ctx context.Context) ([]string, error) {
	var cats []string
	err := r.Conn(ctx).Model(&domain.Product{}).
		Where("category <> ''").
		Distinct().Order("category ASC").Pluck("category", &cats).Error
	if err != nil {
		return nil, apperr.Internal(err, "could not load categories")
	}
	return cats, nil
}
