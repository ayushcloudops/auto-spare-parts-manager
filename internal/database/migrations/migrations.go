// Package migrations holds the ordered, versioned schema migrations applied by
// gormigrate. Each migration has a stable ID and a Rollback, giving real
// upgrade/downgrade paths when a new app version ships to a shop that already
// holds live data — unlike a bare AutoMigrate call.
//
// RULE: never edit a migration that has already shipped. Add a new one.
package migrations

import (
	"autoshop/internal/domain"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// All returns every migration in apply order.
func All() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		initialSchema(),
		addInvoiceItemCostPrice(),
	}
}

// addInvoiceItemCostPrice (0002) adds the cost-price snapshot to invoice items,
// enabling accurate profit reports. Demonstrates the versioned-upgrade path:
// shops already holding data get the new column without losing anything.
func addInvoiceItemCostPrice() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "0002_invoice_item_cost_price",
		Migrate: func(tx *gorm.DB) error {
			return tx.Migrator().AutoMigrate(&domain.InvoiceItem{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&domain.InvoiceItem{}, "cost_price")
		},
	}
}

// initialSchema (0001) creates the complete v1 schema.
func initialSchema() *gormigrate.Migration {
	models := []any{
		&domain.Product{},
		&domain.Customer{},
		&domain.Supplier{},
		&domain.Invoice{},
		&domain.InvoiceItem{},
		&domain.InvoiceSequence{},
		&domain.Purchase{},
		&domain.PurchaseItem{},
		&domain.StockMovement{},
		&domain.ShopProfile{},
		&domain.AppSetting{},
	}
	return &gormigrate.Migration{
		ID: "0001_initial_schema",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(models...)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(models...)
		},
	}
}
