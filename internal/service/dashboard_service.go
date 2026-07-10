package service

import (
	"context"
	"time"

	"autoshop/internal/domain"
)

// DashboardService composes existing repositories to produce the home-screen
// aggregates. It needs no repository of its own — a good sign the data model is
// well-factored.
type DashboardService struct {
	invoices  domain.InvoiceRepository
	products  domain.ProductRepository
	customers domain.CustomerRepository
}

// NewDashboardService wires the service.
func NewDashboardService(invoices domain.InvoiceRepository, products domain.ProductRepository, customers domain.CustomerRepository) *DashboardService {
	return &DashboardService{invoices: invoices, products: products, customers: customers}
}

// Stats returns today's totals plus catalogue/credit figures and recent bills.
func (s *DashboardService) Stats(ctx context.Context) (domain.DashboardStats, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var stats domain.DashboardStats
	var err error

	if stats.TodaySales, err = s.invoices.SumSalesBetween(ctx, startOfDay, endOfDay); err != nil {
		return stats, err
	}
	if stats.TodayBills, err = s.invoices.CountBetween(ctx, startOfDay, endOfDay); err != nil {
		return stats, err
	}
	if stats.TotalProducts, err = s.products.CountAll(ctx); err != nil {
		return stats, err
	}
	if stats.LowStockCount, err = s.products.CountLowStock(ctx); err != nil {
		return stats, err
	}
	if stats.PendingCredit, err = s.customers.SumOutstanding(ctx); err != nil {
		return stats, err
	}
	if stats.RecentInvoices, err = s.invoices.List(ctx, domain.InvoiceFilter{Limit: 5}); err != nil {
		return stats, err
	}
	return stats, nil
}
