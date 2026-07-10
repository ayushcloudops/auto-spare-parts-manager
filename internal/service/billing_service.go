package service

import (
	"context"
	"time"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"
	"autoshop/internal/pkg/money"
)

// BillLineInput is one line the cashier adds to a bill.
type BillLineInput struct {
	ProductID uint        `json:"productId"`
	Quantity  int         `json:"quantity"`
	Discount  money.Money `json:"discount"` // absolute, per line
}

// CreateBillInput is the payload to generate an invoice.
type CreateBillInput struct {
	CustomerID  *uint              `json:"customerId"` // nil = walk-in
	Lines       []BillLineInput    `json:"lines"`
	PaymentMode domain.PaymentMode `json:"paymentMode"`
	AmountPaid  money.Money        `json:"amountPaid"` // for credit/part payment
	Notes       string             `json:"notes"`
}

// BillingService creates invoices. This is the most involved service: a single
// transaction computes GST, allocates the invoice number, writes the invoice
// and its items, decrements stock (with ledger entries) and updates the
// customer's outstanding balance — all atomically.
type BillingService struct {
	invoices  domain.InvoiceRepository
	products  domain.ProductRepository
	customers domain.CustomerRepository
	stock     domain.StockRepository
	settings  domain.SettingsRepository
	tx        Transactor
}

// NewBillingService wires the billing dependencies.
func NewBillingService(
	invoices domain.InvoiceRepository,
	products domain.ProductRepository,
	customers domain.CustomerRepository,
	stock domain.StockRepository,
	settings domain.SettingsRepository,
	tx Transactor,
) *BillingService {
	return &BillingService{
		invoices:  invoices,
		products:  products,
		customers: customers,
		stock:     stock,
		settings:  settings,
		tx:        tx,
	}
}

// isInterState decides IGST vs CGST/SGST by comparing the shop's GST state code
// with the first two digits of the customer's GSTIN. Walk-in or GSTIN-less
// customers are treated as intra-state.
func isInterState(shopStateCode, customerGSTIN string) bool {
	if shopStateCode == "" || len(customerGSTIN) < 2 {
		return false
	}
	return customerGSTIN[:2] != shopStateCode
}

// Create validates and generates an invoice. Returns the saved invoice with its
// items and customer preloaded (ready for display/printing).
func (s *BillingService) Create(ctx context.Context, in CreateBillInput) (*domain.Invoice, error) {
	if len(in.Lines) == 0 {
		return nil, apperr.Validation("add at least one item to the bill")
	}
	if !in.PaymentMode.Valid() {
		return nil, apperr.Validation("select a valid payment mode")
	}

	var invoiceID uint
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		profile, err := s.settings.GetShopProfile(ctx)
		if err != nil {
			return err
		}

		interState := false
		if in.CustomerID != nil {
			cust, err := s.customers.FindByID(ctx, *in.CustomerID)
			if err != nil {
				return err
			}
			interState = isInterState(profile.StateCode, cust.GSTIN)
		}

		// Load products, validate stock, and build lines + snapshots.
		gstLines := make([]GSTLineInput, len(in.Lines))
		items := make([]domain.InvoiceItem, len(in.Lines))
		products := make([]*domain.Product, len(in.Lines))
		for i, l := range in.Lines {
			if l.Quantity <= 0 {
				return apperr.Validation("quantity must be at least 1")
			}
			p, err := s.products.FindByID(ctx, l.ProductID)
			if err != nil {
				return err
			}
			if p.CurrentStock < l.Quantity {
				return apperr.Validation("insufficient stock for %s (only %d left)", p.Name, p.CurrentStock)
			}
			products[i] = p
			gstLines[i] = GSTLineInput{
				UnitPrice: p.SellingPrice,
				Quantity:  l.Quantity,
				Discount:  l.Discount,
				GSTRate:   p.GSTRate,
			}
			pid := p.ID
			items[i] = domain.InvoiceItem{
				ProductID:   &pid,
				ProductName: p.Name,
				PartNumber:  p.PartNumber,
				HSNCode:     p.HSNCode,
				Quantity:    l.Quantity,
				UnitPrice:   p.SellingPrice,
				CostPrice:   p.PurchasePrice, // snapshot for profit reporting
				Discount:    l.Discount,
				GSTRate:     p.GSTRate,
			}
		}

		gst := CalculateGST(gstLines, interState)
		for i := range items {
			lr := gst.Lines[i]
			items[i].TaxableValue = lr.TaxableValue
			items[i].CGST = lr.CGST
			items[i].SGST = lr.SGST
			items[i].IGST = lr.IGST
			items[i].LineTotal = lr.LineTotal
		}

		number, err := s.invoices.NextNumber(ctx, profile.InvoicePrefix, time.Now())
		if err != nil {
			return err
		}

		// Non-credit payments are assumed paid in full unless a partial amount
		// was supplied.
		amountPaid := in.AmountPaid
		if in.PaymentMode != domain.PaymentCredit && amountPaid.IsZero() {
			amountPaid = gst.GrandTotal
		}
		due := gst.GrandTotal.Sub(amountPaid)
		if due.IsNegative() {
			due = money.Zero
		}

		inv := &domain.Invoice{
			Number:        number,
			Date:          time.Now(),
			CustomerID:    in.CustomerID,
			Items:         items,
			PaymentMode:   in.PaymentMode,
			SubTotal:      gst.SubTotal,
			DiscountTotal: gst.DiscountTotal,
			CGST:          gst.CGST,
			SGST:          gst.SGST,
			IGST:          gst.IGST,
			RoundOff:      gst.RoundOff,
			GrandTotal:    gst.GrandTotal,
			AmountPaid:    amountPaid,
			AmountDue:     due,
			Notes:         in.Notes,
		}
		if err := s.invoices.Create(ctx, inv); err != nil {
			return err
		}
		invoiceID = inv.ID

		// Decrement stock and append ledger entries.
		for i, l := range in.Lines {
			p := products[i]
			p.CurrentStock -= l.Quantity
			if err := s.products.Update(ctx, p); err != nil {
				return err
			}
			if err := s.stock.Record(ctx, &domain.StockMovement{
				ProductID:    p.ID,
				Delta:        -l.Quantity,
				BalanceAfter: p.CurrentStock,
				Reason:       domain.StockReasonSale,
				RefType:      "invoice",
				RefID:        inv.ID,
				OccurredAt:   time.Now(),
				Note:         inv.Number,
			}); err != nil {
				return err
			}
		}

		// Credit sale increases the customer's outstanding balance.
		if due.Paise() > 0 && in.CustomerID != nil {
			if err := s.customers.AdjustOutstanding(ctx, *in.CustomerID, due); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.invoices.FindByID(ctx, invoiceID)
}

// Get returns a single invoice with items + customer (for view/reprint).
func (s *BillingService) Get(ctx context.Context, id uint) (*domain.Invoice, error) {
	return s.invoices.FindByID(ctx, id)
}

// List returns invoices matching the filter (history).
func (s *BillingService) List(ctx context.Context, f domain.InvoiceFilter) ([]domain.Invoice, error) {
	return s.invoices.List(ctx, f)
}
