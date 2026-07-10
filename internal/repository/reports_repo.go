package repository

import (
	"context"
	"time"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"
	"autoshop/internal/pkg/money"

	"gorm.io/gorm"
)

type reportsRepo struct {
	db *gorm.DB
}

// NewReportsRepo constructs a ReportsRepository.
func NewReportsRepo(db *gorm.DB) domain.ReportsRepository {
	return &reportsRepo{db: db}
}

func (r *reportsRepo) SalesSummary(ctx context.Context, from, to time.Time) (domain.SalesSummary, error) {
	var row struct {
		Total int64
		Count int64
		Tax   int64
	}
	err := conn(ctx, r.db).Model(&domain.Invoice{}).
		Select("COALESCE(SUM(grand_total),0) AS total, COUNT(*) AS count, COALESCE(SUM(cgst+sgst+igst),0) AS tax").
		Where("date >= ? AND date < ?", from, to).
		Scan(&row).Error
	if err != nil {
		return domain.SalesSummary{}, apperr.Internal(err, "sales summary failed")
	}
	return domain.SalesSummary{
		From:         from,
		To:           to,
		TotalSales:   money.FromPaise(row.Total),
		InvoiceCount: row.Count,
		TotalTax:     money.FromPaise(row.Tax),
	}, nil
}

func (r *reportsRepo) TopProducts(ctx context.Context, from, to time.Time, limit int) ([]domain.TopProduct, error) {
	if limit <= 0 {
		limit = 10
	}
	var rows []struct {
		ProductID   uint
		ProductName string
		QtySold     int
		Revenue     int64
	}
	err := conn(ctx, r.db).Table("invoice_items AS ii").
		Select("ii.product_id AS product_id, ii.product_name AS product_name, SUM(ii.quantity) AS qty_sold, SUM(ii.line_total) AS revenue").
		Joins("JOIN invoices i ON i.id = ii.invoice_id").
		Where("i.date >= ? AND i.date < ?", from, to).
		Where("i.deleted_at IS NULL AND ii.deleted_at IS NULL").
		Group("ii.product_id, ii.product_name").
		Order("qty_sold DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, apperr.Internal(err, "top products failed")
	}
	out := make([]domain.TopProduct, len(rows))
	for i, r := range rows {
		out[i] = domain.TopProduct{
			ProductID:   r.ProductID,
			ProductName: r.ProductName,
			QtySold:     r.QtySold,
			Revenue:     money.FromPaise(r.Revenue),
		}
	}
	return out, nil
}

func (r *reportsRepo) Profit(ctx context.Context, from, to time.Time) (domain.ProfitReport, error) {
	var row struct {
		Revenue int64
		Cost    int64
	}
	err := conn(ctx, r.db).Table("invoice_items AS ii").
		Select("COALESCE(SUM(ii.taxable_value),0) AS revenue, COALESCE(SUM(ii.cost_price * ii.quantity),0) AS cost").
		Joins("JOIN invoices i ON i.id = ii.invoice_id").
		Where("i.date >= ? AND i.date < ?", from, to).
		Where("i.deleted_at IS NULL AND ii.deleted_at IS NULL").
		Scan(&row).Error
	if err != nil {
		return domain.ProfitReport{}, apperr.Internal(err, "profit report failed")
	}
	return domain.ProfitReport{
		From:    from,
		To:      to,
		Revenue: money.FromPaise(row.Revenue),
		Cost:    money.FromPaise(row.Cost),
		Profit:  money.FromPaise(row.Revenue - row.Cost),
	}, nil
}
