package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"time"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"
)

// ReportsService provides sales/profit/inventory reporting and CSV export.
type ReportsService struct {
	reports  domain.ReportsRepository
	products domain.ProductRepository
	invoices domain.InvoiceRepository
}

// NewReportsService wires the service.
func NewReportsService(reports domain.ReportsRepository, products domain.ProductRepository, invoices domain.InvoiceRepository) *ReportsService {
	return &ReportsService{reports: reports, products: products, invoices: invoices}
}

// Sales returns the sales summary for [from, to).
func (s *ReportsService) Sales(ctx context.Context, from, to time.Time) (domain.SalesSummary, error) {
	return s.reports.SalesSummary(ctx, from, to)
}

// TopProducts returns the best-selling products for the period.
func (s *ReportsService) TopProducts(ctx context.Context, from, to time.Time, limit int) ([]domain.TopProduct, error) {
	return s.reports.TopProducts(ctx, from, to, limit)
}

// Profit returns the profit report for the period.
func (s *ReportsService) Profit(ctx context.Context, from, to time.Time) (domain.ProfitReport, error) {
	return s.reports.Profit(ctx, from, to)
}

// LowStock returns products at or below their minimum stock.
func (s *ReportsService) LowStock(ctx context.Context) ([]domain.Product, error) {
	return s.products.List(ctx, domain.ProductFilter{LowStockOnly: true})
}

// ExportSalesCSV builds a CSV of invoices in [from, to). Returned as a string
// the frontend saves to disk.
func (s *ReportsService) ExportSalesCSV(ctx context.Context, from, to time.Time) (string, error) {
	invoices, err := s.invoices.List(ctx, domain.InvoiceFilter{From: &from, To: &to})
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"Invoice", "Date", "Customer", "Payment", "Taxable", "CGST", "SGST", "IGST", "Grand Total", "Paid", "Due"})
	for _, inv := range invoices {
		customer := "Walk-in"
		if inv.Customer != nil {
			customer = inv.Customer.Name
		}
		_ = w.Write([]string{
			inv.Number,
			inv.Date.Format("2006-01-02 15:04"),
			customer,
			string(inv.PaymentMode),
			inv.SubTotal.String(),
			inv.CGST.String(),
			inv.SGST.String(),
			inv.IGST.String(),
			inv.GrandTotal.String(),
			inv.AmountPaid.String(),
			inv.AmountDue.String(),
		})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", apperr.Internal(err, "csv export failed")
	}
	return buf.String(), nil
}
