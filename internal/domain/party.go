package domain

import (
	"context"

	"autoshop/internal/pkg/money"
)

// Customer is a buyer. Outstanding tracks unpaid credit; CreditLimit caps how
// much credit the shop will extend. Both are stored in paise.
type Customer struct {
	Base
	Name        string      `gorm:"index;not null" json:"name"`
	Phone       string      `gorm:"index" json:"phone"`
	Address     string      `json:"address"`
	GSTIN       string      `json:"gstin"` // customer's GST number (B2B)
	Outstanding money.Money `gorm:"not null;default:0" json:"outstanding"`
	CreditLimit money.Money `gorm:"not null;default:0" json:"creditLimit"`
}

// Supplier is a vendor the shop buys stock from.
type Supplier struct {
	Base
	Name    string `gorm:"index;not null" json:"name"`
	Phone   string `gorm:"index" json:"phone"`
	Address string `json:"address"`
	GSTIN   string `json:"gstin"`
}

// CustomerRepository persists customers and their outstanding credit balances.
type CustomerRepository interface {
	Create(ctx context.Context, c *Customer) error
	Update(ctx context.Context, c *Customer) error
	FindByID(ctx context.Context, id uint) (*Customer, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, search string) ([]Customer, error)
	// AdjustOutstanding atomically adds delta (may be negative) to a customer's
	// outstanding balance — used when a bill is on credit or a payment is made.
	AdjustOutstanding(ctx context.Context, id uint, delta money.Money) error
	SumOutstanding(ctx context.Context) (money.Money, error)
}

// SupplierRepository persists suppliers.
type SupplierRepository interface {
	Create(ctx context.Context, s *Supplier) error
	Update(ctx context.Context, s *Supplier) error
	FindByID(ctx context.Context, id uint) (*Supplier, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, search string) ([]Supplier, error)
}
