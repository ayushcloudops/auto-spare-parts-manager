package app

import (
	"context"
	"time"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"
	"autoshop/internal/service"
)

// ReportsHandler is the Wails binding for the Reports module. Date parameters
// are "YYYY-MM-DD" strings (from the date pickers); the range is treated as the
// whole days from..to inclusive.
type ReportsHandler struct {
	svc *service.ReportsService
}

// NewReportsHandler constructs the handler.
func NewReportsHandler(svc *service.ReportsService) *ReportsHandler {
	return &ReportsHandler{svc: svc}
}

func (h *ReportsHandler) ctx() context.Context { return context.Background() }

// parseRange turns two YYYY-MM-DD strings into a [from, to) time window where to
// is exclusive end-of-day.
func parseRange(fromStr, toStr string) (time.Time, time.Time, error) {
	from, err := time.ParseInLocation("2006-01-02", fromStr, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, apperr.Validation("invalid from date")
	}
	to, err := time.ParseInLocation("2006-01-02", toStr, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, apperr.Validation("invalid to date")
	}
	return from, to.Add(24 * time.Hour), nil // include the whole 'to' day
}

// Sales returns the sales summary for the date range.
func (h *ReportsHandler) Sales(from, to string) (domain.SalesSummary, error) {
	f, t, err := parseRange(from, to)
	if err != nil {
		return domain.SalesSummary{}, bindError(err)
	}
	s, err := h.svc.Sales(h.ctx(), f, t)
	return s, bindError(err)
}

// TopProducts returns the best sellers for the range.
func (h *ReportsHandler) TopProducts(from, to string, limit int) ([]domain.TopProduct, error) {
	f, t, err := parseRange(from, to)
	if err != nil {
		return nil, bindError(err)
	}
	items, err := h.svc.TopProducts(h.ctx(), f, t, limit)
	return items, bindError(err)
}

// Profit returns the profit report for the range.
func (h *ReportsHandler) Profit(from, to string) (domain.ProfitReport, error) {
	f, t, err := parseRange(from, to)
	if err != nil {
		return domain.ProfitReport{}, bindError(err)
	}
	p, err := h.svc.Profit(h.ctx(), f, t)
	return p, bindError(err)
}

// LowStock returns products at or below minimum stock.
func (h *ReportsHandler) LowStock() ([]domain.Product, error) {
	items, err := h.svc.LowStock(h.ctx())
	return items, bindError(err)
}

// ExportSalesCSV returns a CSV string of invoices in the range.
func (h *ReportsHandler) ExportSalesCSV(from, to string) (string, error) {
	f, t, err := parseRange(from, to)
	if err != nil {
		return "", bindError(err)
	}
	csv, err := h.svc.ExportSalesCSV(h.ctx(), f, t)
	return csv, bindError(err)
}
