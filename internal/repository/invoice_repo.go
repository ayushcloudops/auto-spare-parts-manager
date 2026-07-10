package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"
	"autoshop/internal/pkg/money"

	"gorm.io/gorm"
)

type invoiceRepo struct {
	Base[domain.Invoice]
}

// NewInvoiceRepo constructs an InvoiceRepository.
func NewInvoiceRepo(db *gorm.DB) domain.InvoiceRepository {
	return &invoiceRepo{Base: NewBase[domain.Invoice](db)}
}

// financialYear returns the Indian FY code for t, e.g. "2526" for FY 2025-26.
// The Indian financial year runs 1 April – 31 March.
func financialYear(t time.Time) string {
	y := t.Year()
	if int(t.Month()) < int(time.April) {
		y-- // Jan–Mar fall in the previous financial year
	}
	return fmt.Sprintf("%02d%02d", y%100, (y+1)%100)
}

func (r *invoiceRepo) FindByID(ctx context.Context, id uint) (*domain.Invoice, error) {
	var inv domain.Invoice
	err := r.Conn(ctx).Preload("Items").Preload("Customer").First(&inv, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.NotFound("invoice #%d not found", id)
	}
	if err != nil {
		return nil, apperr.Internal(err, "could not load invoice")
	}
	return &inv, nil
}

func (r *invoiceRepo) List(ctx context.Context, f domain.InvoiceFilter) ([]domain.Invoice, error) {
	q := r.Conn(ctx).Model(&domain.Invoice{}).
		Preload("Customer").
		Joins("LEFT JOIN customers ON customers.id = invoices.customer_id").
		Order("invoices.date DESC, invoices.id DESC")

	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("invoices.number LIKE ? OR customers.name LIKE ?", like, like)
	}
	if f.CustomerID != nil {
		q = q.Where("invoices.customer_id = ?", *f.CustomerID)
	}
	if f.From != nil {
		q = q.Where("invoices.date >= ?", *f.From)
	}
	if f.To != nil {
		q = q.Where("invoices.date <= ?", *f.To)
	}
	if f.Limit > 0 {
		q = q.Limit(f.Limit).Offset(f.Offset)
	}

	var out []domain.Invoice
	if err := q.Find(&out).Error; err != nil {
		return nil, apperr.Internal(err, "could not list invoices")
	}
	return out, nil
}

// NextNumber allocates the next invoice number for the current financial year.
// It must be called inside the same transaction as the invoice insert so the
// increment and the invoice are committed together — guaranteeing uniqueness.
func (r *invoiceRepo) NextNumber(ctx context.Context, prefix string, now time.Time) (string, error) {
	fy := financialYear(now)
	tx := r.Conn(ctx)

	var seq domain.InvoiceSequence
	if err := tx.Where(domain.InvoiceSequence{FY: fy}).FirstOrCreate(&seq).Error; err != nil {
		return "", apperr.Internal(err, "could not read invoice sequence")
	}
	seq.LastNumber++
	if err := tx.Save(&seq).Error; err != nil {
		return "", apperr.Internal(err, "could not update invoice sequence")
	}
	if prefix == "" {
		prefix = "INV"
	}
	return fmt.Sprintf("%s-%s-%04d", prefix, fy, seq.LastNumber), nil
}

func (r *invoiceRepo) CountBetween(ctx context.Context, from, to time.Time) (int64, error) {
	var n int64
	err := r.Conn(ctx).Model(&domain.Invoice{}).
		Where("date >= ? AND date < ?", from, to).Count(&n).Error
	if err != nil {
		return 0, apperr.Internal(err, "could not count invoices")
	}
	return n, nil
}

func (r *invoiceRepo) SumSalesBetween(ctx context.Context, from, to time.Time) (money.Money, error) {
	var total int64
	err := r.Conn(ctx).Model(&domain.Invoice{}).
		Where("date >= ? AND date < ?", from, to).
		Select("COALESCE(SUM(grand_total), 0)").Scan(&total).Error
	if err != nil {
		return 0, apperr.Internal(err, "could not sum sales")
	}
	return money.FromPaise(total), nil
}
