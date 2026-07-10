package app

import (
	"context"

	"autoshop/internal/domain"
	"autoshop/internal/service"
)

// BillingHandler is the Wails binding for creating and reading invoices. Billing
// and Invoice History both use it.
type BillingHandler struct {
	svc *service.BillingService
}

// NewBillingHandler constructs the handler.
func NewBillingHandler(svc *service.BillingService) *BillingHandler {
	return &BillingHandler{svc: svc}
}

func (h *BillingHandler) ctx() context.Context { return context.Background() }

// CreateBill generates a new invoice from the cart and returns it (with items).
func (h *BillingHandler) CreateBill(in service.CreateBillInput) (*domain.Invoice, error) {
	inv, err := h.svc.Create(h.ctx(), in)
	return inv, bindError(err)
}

// GetInvoice returns one invoice by id (for view/reprint).
func (h *BillingHandler) GetInvoice(id uint) (*domain.Invoice, error) {
	inv, err := h.svc.Get(h.ctx(), id)
	return inv, bindError(err)
}

// ListInvoices returns invoices matching the filter (history).
func (h *BillingHandler) ListInvoices(filter domain.InvoiceFilter) ([]domain.Invoice, error) {
	items, err := h.svc.List(h.ctx(), filter)
	return items, bindError(err)
}
