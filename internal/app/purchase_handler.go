package app

import (
	"context"

	"autoshop/internal/domain"
	"autoshop/internal/service"
)

// PurchaseHandler is the Wails binding for the Purchase Entry module.
type PurchaseHandler struct {
	svc *service.PurchaseService
}

// NewPurchaseHandler constructs the handler.
func NewPurchaseHandler(svc *service.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{svc: svc}
}

func (h *PurchaseHandler) ctx() context.Context { return context.Background() }

// Create records a purchase (increases stock).
func (h *PurchaseHandler) Create(in service.CreatePurchaseInput) (*domain.Purchase, error) {
	p, err := h.svc.Create(h.ctx(), in)
	return p, bindError(err)
}

// Get returns a purchase with items + supplier.
func (h *PurchaseHandler) Get(id uint) (*domain.Purchase, error) {
	p, err := h.svc.Get(h.ctx(), id)
	return p, bindError(err)
}

// List returns purchases matching the filter.
func (h *PurchaseHandler) List(filter domain.PurchaseFilter) ([]domain.Purchase, error) {
	items, err := h.svc.List(h.ctx(), filter)
	return items, bindError(err)
}
