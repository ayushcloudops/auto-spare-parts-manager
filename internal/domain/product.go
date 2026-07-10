package domain

import (
	"context"

	"autoshop/internal/pkg/money"
)

// Product is a spare part held in inventory.
//
// Prices are stored as money.Money (paise). GSTRate is the GST percentage that
// applies to this part (one of the standard slabs: 0, 5, 12, 18, 28). HSNCode
// is the Harmonised System of Nomenclature code required on GST tax invoices.
type Product struct {
	Base
	Name          string  `gorm:"index;not null" json:"name"`
	PartNumber    string  `gorm:"index" json:"partNumber"`
	Brand         string  `json:"brand"`
	VehicleBrand  string  `gorm:"index" json:"vehicleBrand"`
	VehicleModel  string  `json:"vehicleModel"`
	VehicleYear   string  `json:"vehicleYear"` // free text, e.g. "2015-2020"
	Category      string  `gorm:"index" json:"category"`
	HSNCode       string  `json:"hsnCode"`
	PurchasePrice money.Money `json:"purchasePrice"`
	SellingPrice  money.Money `json:"sellingPrice"`
	GSTRate       float64 `gorm:"not null;default:0" json:"gstRate"`
	CurrentStock  int     `gorm:"not null;default:0" json:"currentStock"`
	MinimumStock  int     `gorm:"not null;default:0" json:"minimumStock"`
	Location      string  `json:"location"` // rack / shelf
}

// IsLowStock reports whether current stock is at or below the minimum.
func (p Product) IsLowStock() bool {
	return p.CurrentStock <= p.MinimumStock
}

// ProductFilter describes the criteria for listing/searching products.
// Zero-value fields are ignored, so the same struct serves "list all",
// "search by text" and "show only low stock".
type ProductFilter struct {
	Search       string `json:"search"`       // matches name / part number / brand
	Category     string `json:"category"`     // exact category
	VehicleBrand string `json:"vehicleBrand"` // exact vehicle brand
	LowStockOnly bool   `json:"lowStockOnly"` // current <= minimum
	Limit        int    `json:"limit"`        // page size (0 = no limit)
	Offset       int    `json:"offset"`       // page offset
}

// ProductRepository describes persistence for products. The service layer
// depends on this interface, not on the GORM implementation, so storage can be
// swapped (e.g. a future cloud backend) without touching business logic.
type ProductRepository interface {
	Create(ctx context.Context, p *Product) error
	Update(ctx context.Context, p *Product) error
	FindByID(ctx context.Context, id uint) (*Product, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, f ProductFilter) ([]Product, error)
	Count(ctx context.Context, f ProductFilter) (int64, error)
	CountAll(ctx context.Context) (int64, error)
	CountLowStock(ctx context.Context) (int64, error)
	Categories(ctx context.Context) ([]string, error)
}
