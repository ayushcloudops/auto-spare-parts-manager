// Package domain holds the core business entities and the repository
// interfaces that describe how they are persisted. It deliberately contains NO
// database driver code beyond GORM struct tags — the concrete persistence logic
// lives in internal/repository. This keeps the domain the stable centre of the
// application that everything else depends on (Clean Architecture).
package domain

import (
	"time"

	"gorm.io/gorm"
)

// Base is embedded in every persisted entity. It provides an auto-increment
// primary key, audit timestamps and soft-delete support (DeletedAt). Soft
// deletes matter for a shop: a "deleted" product that already appears on past
// invoices must not vanish from history.
type Base struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// PaymentMode enumerates how a bill was settled.
type PaymentMode string

const (
	PaymentCash   PaymentMode = "cash"
	PaymentUPI    PaymentMode = "upi"
	PaymentCard   PaymentMode = "card"
	PaymentCredit PaymentMode = "credit" // added to customer outstanding
)

// Valid reports whether the payment mode is one of the known values.
func (p PaymentMode) Valid() bool {
	switch p {
	case PaymentCash, PaymentUPI, PaymentCard, PaymentCredit:
		return true
	default:
		return false
	}
}

// StockReason explains why a StockMovement row was written.
type StockReason string

const (
	StockReasonSale       StockReason = "sale"
	StockReasonPurchase   StockReason = "purchase"
	StockReasonAdjustment StockReason = "adjustment"
	StockReasonReturn     StockReason = "return"
	StockReasonOpening    StockReason = "opening" // initial stock when product created
)
