package domain

import (
	"context"
	"time"

	"autoshop/internal/pkg/money"
)

// SalesSummary aggregates sales over a period (daily/weekly/monthly reports).
type SalesSummary struct {
	From         time.Time   `json:"from"`
	To           time.Time   `json:"to"`
	TotalSales   money.Money `json:"totalSales"`
	InvoiceCount int64       `json:"invoiceCount"`
	TotalTax     money.Money `json:"totalTax"`
}

// TopProduct is a row in the best-sellers report.
type TopProduct struct {
	ProductID   uint        `json:"productId"`
	ProductName string      `json:"productName"`
	QtySold     int         `json:"qtySold"`
	Revenue     money.Money `json:"revenue"`
}

// ProfitReport summarises revenue vs cost of goods sold over a period.
type ProfitReport struct {
	From    time.Time   `json:"from"`
	To      time.Time   `json:"to"`
	Revenue money.Money `json:"revenue"`
	Cost    money.Money `json:"cost"`
	Profit  money.Money `json:"profit"`
}

// ReportsRepository runs the cross-entity aggregate queries used by reporting.
type ReportsRepository interface {
	SalesSummary(ctx context.Context, from, to time.Time) (SalesSummary, error)
	TopProducts(ctx context.Context, from, to time.Time, limit int) ([]TopProduct, error)
	Profit(ctx context.Context, from, to time.Time) (ProfitReport, error)
}

// DashboardStats is the aggregate shown on the home dashboard.
type DashboardStats struct {
	TodaySales     money.Money `json:"todaySales"`
	TodayBills     int64       `json:"todayBills"`
	TotalProducts  int64       `json:"totalProducts"`
	LowStockCount  int64       `json:"lowStockCount"`
	PendingCredit  money.Money `json:"pendingCredit"`
	RecentInvoices []Invoice   `json:"recentInvoices"`
}
