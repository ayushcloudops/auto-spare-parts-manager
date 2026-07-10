package domain

import (
	"context"
	"time"
)

// StockMovement is an append-only ledger of every change to a product's stock.
// Each sale, purchase, return or manual adjustment writes one row. This makes
// the current stock figure auditable and reconstructable — essential when a
// shopkeeper disputes a count. BalanceAfter records the resulting on-hand
// quantity for quick verification.
type StockMovement struct {
	Base
	ProductID    uint        `gorm:"index;not null" json:"productId"`
	Delta        int         `gorm:"not null" json:"delta"` // +received / -sold
	BalanceAfter int         `gorm:"not null" json:"balanceAfter"`
	Reason       StockReason `gorm:"not null" json:"reason"`
	RefType      string      `json:"refType"` // "invoice" | "purchase" | ""
	RefID        uint        `json:"refId"`   // id of the referenced document
	Note         string      `json:"note"`
	OccurredAt   time.Time   `gorm:"index;not null" json:"occurredAt"`
}

// StockRepository persists the stock ledger. Sales, purchases and adjustments
// all append movements through it.
type StockRepository interface {
	Record(ctx context.Context, m *StockMovement) error
	ListByProduct(ctx context.Context, productID uint) ([]StockMovement, error)
}
