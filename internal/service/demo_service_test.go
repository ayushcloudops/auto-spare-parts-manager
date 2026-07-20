package service_test

import (
	"context"
	"path/filepath"
	"testing"

	"autoshop/internal/database"
	"autoshop/internal/domain"
	"autoshop/internal/repository"
	"autoshop/internal/service"
)

// TestDemoDataLoads verifies the sample dataset lands and is internally
// consistent — because it is created through the real services, stock, GST and
// customer credit must all agree with the invoices generated.
func TestDemoDataLoads(t *testing.T) {
	ctx := context.Background()
	db, err := database.Open(filepath.Join(t.TempDir(), "demo.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := database.Seed(db); err != nil {
		t.Fatalf("seed: %v", err)
	}

	tx := repository.NewTxManager(db)
	productRepo := repository.NewProductRepo(db)
	stockRepo := repository.NewStockRepo(db)
	invoiceRepo := repository.NewInvoiceRepo(db)
	customerRepo := repository.NewCustomerRepo(db)
	supplierRepo := repository.NewSupplierRepo(db)
	settingsRepo := repository.NewSettingsRepo(db)

	productSvc := service.NewProductService(productRepo, stockRepo, tx)
	billingSvc := service.NewBillingService(invoiceRepo, productRepo, customerRepo, stockRepo, settingsRepo, tx)
	customerSvc := service.NewCustomerService(customerRepo, invoiceRepo)
	supplierSvc := service.NewSupplierService(supplierRepo)
	settingsSvc := service.NewSettingsService(settingsRepo)

	demo := service.NewDemoService(productSvc, customerSvc, supplierSvc, billingSvc, settingsSvc)

	summary, err := demo.Load(ctx)
	if err != nil {
		t.Fatalf("load demo data: %v", err)
	}

	if summary.Products != 12 || summary.Customers != 4 || summary.Suppliers != 3 || summary.Invoices != 4 {
		t.Fatalf("unexpected summary: %+v", summary)
	}

	// Shop profile should have been personalised for the demo.
	profile, err := settingsSvc.GetShopProfile(ctx)
	if err != nil || profile.StateCode != "27" {
		t.Fatalf("shop profile not set for demo: %+v (err %v)", profile, err)
	}

	// Low-stock items should exist so the alerts/report have content.
	low, err := productRepo.CountLowStock(ctx)
	if err != nil {
		t.Fatalf("count low stock: %v", err)
	}
	if low < 2 {
		t.Fatalf("expected at least 2 low-stock products, got %d", low)
	}

	// The credit sale must have left an outstanding balance.
	outstanding, err := customerRepo.SumOutstanding(ctx)
	if err != nil {
		t.Fatalf("sum outstanding: %v", err)
	}
	if outstanding.Paise() <= 0 {
		t.Fatalf("expected outstanding credit from the demo credit sale, got %s", outstanding)
	}

	// Stock must reflect the sales: oil filter started at 40, 2 were sold.
	items, err := productRepo.List(ctx, domain.ProductFilter{Search: "Oil Filter"})
	if err != nil || len(items) == 0 {
		t.Fatalf("oil filter missing: %v", err)
	}
	if items[0].CurrentStock != 38 {
		t.Fatalf("oil filter stock should be 38 after selling 2, got %d", items[0].CurrentStock)
	}
}
