// Package app is the composition root. It wires concrete implementations
// together (dependency injection) and will hold the Wails binding handlers —
// one per module — that the frontend calls.
//
// Keeping all wiring here means there is exactly one place that knows how the
// application is assembled; the rest of the code depends only on interfaces.
package app

import (
	"autoshop/internal/config"
	"autoshop/internal/database"
	"autoshop/internal/printer"
	"autoshop/internal/repository"
	"autoshop/internal/service"

	"gorm.io/gorm"
)

// Container holds shared infrastructure and the module services/handlers. As
// each module is built, its handler is added here and exposed to Wails via
// Handlers().
type Container struct {
	DB *gorm.DB
	Tx *repository.TxManager

	// System is the app-level health/identity handler.
	System *SystemHandler

	// Products is the Wails handler for the product/inventory module.
	Products *ProductHandler

	// Billing handles invoice creation + history; Customers handles customers.
	Billing   *BillingHandler
	Customers *CustomerHandler

	// Remaining module handlers.
	Suppliers *SupplierHandler
	Purchases *PurchaseHandler
	Reports   *ReportsHandler
	Dashboard *DashboardHandler
	Settings  *SettingsHandler
	Print     *PrintHandler
	Demo      *DemoHandler
}

// Bootstrap opens the database, applies migrations and seed data, and assembles
// the dependency graph. Called once at startup before the UI launches.
func Bootstrap() (*Container, error) {
	dbPath, err := config.DBPath()
	if err != nil {
		return nil, err
	}

	db, err := database.Open(dbPath)
	if err != nil {
		return nil, err
	}
	if err := database.Migrate(db); err != nil {
		return nil, err
	}
	if err := database.Seed(db); err != nil {
		return nil, err
	}

	tx := repository.NewTxManager(db)

	// --- wire repositories -> services -> handlers (dependency injection) ---
	productRepo := repository.NewProductRepo(db)
	stockRepo := repository.NewStockRepo(db)
	invoiceRepo := repository.NewInvoiceRepo(db)
	customerRepo := repository.NewCustomerRepo(db)
	supplierRepo := repository.NewSupplierRepo(db)
	purchaseRepo := repository.NewPurchaseRepo(db)
	settingsRepo := repository.NewSettingsRepo(db)
	reportsRepo := repository.NewReportsRepo(db)

	productSvc := service.NewProductService(productRepo, stockRepo, tx)
	billingSvc := service.NewBillingService(invoiceRepo, productRepo, customerRepo, stockRepo, settingsRepo, tx)
	customerSvc := service.NewCustomerService(customerRepo, invoiceRepo)
	supplierSvc := service.NewSupplierService(supplierRepo)
	purchaseSvc := service.NewPurchaseService(purchaseRepo, productRepo, stockRepo, tx)
	settingsSvc := service.NewSettingsService(settingsRepo)
	reportsSvc := service.NewReportsService(reportsRepo, productRepo, invoiceRepo)
	dashboardSvc := service.NewDashboardService(invoiceRepo, productRepo, customerRepo)
	demoSvc := service.NewDemoService(productSvc, customerSvc, supplierSvc, billingSvc, settingsSvc)

	c := &Container{
		DB:        db,
		Tx:        tx,
		System:    NewSystemHandler(db),
		Products:  NewProductHandler(productSvc),
		Billing:   NewBillingHandler(billingSvc),
		Customers: NewCustomerHandler(customerSvc),
		Suppliers: NewSupplierHandler(supplierSvc),
		Purchases: NewPurchaseHandler(purchaseSvc),
		Reports:   NewReportsHandler(reportsSvc),
		Dashboard: NewDashboardHandler(dashboardSvc),
		Settings:  NewSettingsHandler(settingsSvc),
		Print:     NewPrintHandler(invoiceRepo, settingsRepo, printer.New()),
		Demo:      NewDemoHandler(demoSvc),
	}

	return c, nil
}

// Handlers returns the list of structs Wails should bind and expose to the
// frontend. Currently empty; each module appends its handler. The root App
// (in package main) is bound separately.
func (c *Container) Handlers() []interface{} {
	return []interface{}{
		c.System,
		c.Products,
		c.Billing,
		c.Customers,
		c.Suppliers,
		c.Purchases,
		c.Reports,
		c.Dashboard,
		c.Settings,
		c.Print,
		c.Demo,
	}
}

// Close releases the database connection. Called on application shutdown.
func (c *Container) Close() error {
	if c.DB == nil {
		return nil
	}
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
