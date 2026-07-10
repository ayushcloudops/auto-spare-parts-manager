package domain

import (
	"context"
	"time"

	"autoshop/internal/pkg/money"
)

// Purchase records a supplier invoice (goods received). Saving a purchase
// increases stock for each line and updates the product's purchase price.
type Purchase struct {
	Base
	SupplierID    uint           `gorm:"index;not null" json:"supplierId"`
	Supplier      *Supplier      `json:"supplier,omitempty"`
	SupplierInvNo string         `json:"supplierInvNo"` // the supplier's own bill number
	Date          time.Time      `gorm:"index;not null" json:"date"`
	Items         []PurchaseItem `gorm:"constraint:OnDelete:CASCADE" json:"items"`

	SubTotal   money.Money `gorm:"not null;default:0" json:"subTotal"`
	GSTTotal   money.Money `gorm:"not null;default:0" json:"gstTotal"`
	GrandTotal money.Money `gorm:"not null;default:0" json:"grandTotal"`
	Notes      string      `json:"notes"`
}

// PurchaseItem is one received line. CostPrice is what the shop paid per unit.
type PurchaseItem struct {
	Base
	PurchaseID  uint        `gorm:"index;not null" json:"purchaseId"`
	ProductID   uint        `gorm:"index;not null" json:"productId"`
	ProductName string      `gorm:"not null" json:"productName"` // snapshot
	Quantity    int         `gorm:"not null" json:"quantity"`
	CostPrice   money.Money `gorm:"not null" json:"costPrice"`
	GSTRate     float64     `gorm:"not null;default:0" json:"gstRate"`
	LineTotal   money.Money `gorm:"not null" json:"lineTotal"`
}

// PurchaseFilter describes criteria for listing purchases.
type PurchaseFilter struct {
	SupplierID *uint `json:"supplierId"`
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
}

// PurchaseRepository persists supplier purchases (goods received).
type PurchaseRepository interface {
	Create(ctx context.Context, p *Purchase) error
	FindByID(ctx context.Context, id uint) (*Purchase, error) // preloads items + supplier
	List(ctx context.Context, f PurchaseFilter) ([]Purchase, error)
}
