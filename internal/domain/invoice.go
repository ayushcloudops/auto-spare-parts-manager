package domain

import (
	"context"
	"time"

	"autoshop/internal/pkg/money"
)

// Invoice is a sales bill. Monetary fields are denormalised totals computed by
// the billing service so that history is immutable: even if a product's price
// or GST rate changes later, a past invoice still shows what was actually
// charged.
//
// GST split: for an intra-state sale the tax is divided equally into CGST and
// SGST; for an inter-state sale it is charged as IGST. The billing service
// decides which based on the shop's state vs the customer's state.
type Invoice struct {
	Base
	Number      string      `gorm:"uniqueIndex;not null" json:"number"`
	Date        time.Time   `gorm:"index;not null" json:"date"`
	CustomerID  *uint       `gorm:"index" json:"customerId"` // nil = walk-in customer
	Customer    *Customer   `json:"customer,omitempty"`
	Items       []InvoiceItem `gorm:"constraint:OnDelete:CASCADE" json:"items"`
	PaymentMode PaymentMode `gorm:"not null" json:"paymentMode"`

	SubTotal      money.Money `gorm:"not null;default:0" json:"subTotal"`      // taxable value (after line discounts, before tax)
	DiscountTotal money.Money `gorm:"not null;default:0" json:"discountTotal"` // sum of all discounts
	CGST          money.Money `gorm:"not null;default:0" json:"cgst"`
	SGST          money.Money `gorm:"not null;default:0" json:"sgst"`
	IGST          money.Money `gorm:"not null;default:0" json:"igst"`
	RoundOff      money.Money `gorm:"not null;default:0" json:"roundOff"` // to a whole rupee grand total
	GrandTotal    money.Money `gorm:"not null;default:0" json:"grandTotal"`
	AmountPaid    money.Money `gorm:"not null;default:0" json:"amountPaid"`
	AmountDue     money.Money `gorm:"not null;default:0" json:"amountDue"` // unpaid portion (credit)

	Notes string `json:"notes"`
}

// InvoiceItem is a single line on an invoice. ProductName/PartNumber/HSNCode are
// snapshots taken at sale time so deleting/editing the product never alters the
// historical bill.
type InvoiceItem struct {
	Base
	InvoiceID   uint        `gorm:"index;not null" json:"invoiceId"`
	ProductID   *uint       `gorm:"index" json:"productId"` // nil if product later hard-deleted
	ProductName string      `gorm:"not null" json:"productName"`
	PartNumber  string      `json:"partNumber"`
	HSNCode     string      `json:"hsnCode"`
	Quantity    int         `gorm:"not null" json:"quantity"`
	UnitPrice   money.Money `gorm:"not null" json:"unitPrice"`
	CostPrice   money.Money `gorm:"not null;default:0" json:"costPrice"` // snapshot of purchase price, for profit reports
	Discount    money.Money `gorm:"not null;default:0" json:"discount"`  // absolute, on the line
	GSTRate     float64     `gorm:"not null;default:0" json:"gstRate"`
	TaxableValue money.Money `gorm:"not null" json:"taxableValue"` // qty*price - discount
	CGST        money.Money `gorm:"not null;default:0" json:"cgst"`
	SGST        money.Money `gorm:"not null;default:0" json:"sgst"`
	IGST        money.Money `gorm:"not null;default:0" json:"igst"`
	LineTotal   money.Money `gorm:"not null" json:"lineTotal"` // taxable + taxes
}

// InvoiceSequence provides gap-free, per-financial-year invoice numbering.
// The billing service increments LastNumber inside the same transaction that
// writes the invoice, guaranteeing uniqueness even under rapid billing.
// FY is the Indian financial year code, e.g. "2526" for FY 2025-26.
type InvoiceSequence struct {
	FY         string `gorm:"primaryKey" json:"fy"`
	LastNumber int    `gorm:"not null;default:0" json:"lastNumber"`
}

// InvoiceFilter describes criteria for listing invoices (history search).
type InvoiceFilter struct {
	Search     string     `json:"search"` // invoice number or customer name
	CustomerID *uint      `json:"customerId"`
	From       *time.Time `json:"from"`
	To         *time.Time `json:"to"`
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset"`
}

// InvoiceRepository persists sales invoices and issues invoice numbers.
type InvoiceRepository interface {
	Create(ctx context.Context, inv *Invoice) error
	FindByID(ctx context.Context, id uint) (*Invoice, error) // preloads items + customer
	List(ctx context.Context, f InvoiceFilter) ([]Invoice, error)
	// NextNumber atomically allocates the next per-financial-year invoice number.
	NextNumber(ctx context.Context, prefix string, now time.Time) (string, error)
	CountBetween(ctx context.Context, from, to time.Time) (int64, error)
	SumSalesBetween(ctx context.Context, from, to time.Time) (money.Money, error)
}
