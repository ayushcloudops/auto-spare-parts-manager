package app

import (
	"context"

	"autoshop/internal/domain"
	"autoshop/internal/service"
)

// SupplierHandler is the Wails binding for the Suppliers module.
type SupplierHandler struct {
	svc *service.SupplierService
}

// NewSupplierHandler constructs the handler.
func NewSupplierHandler(svc *service.SupplierService) *SupplierHandler {
	return &SupplierHandler{svc: svc}
}

func (h *SupplierHandler) ctx() context.Context { return context.Background() }

func (h *SupplierHandler) List(search string) ([]domain.Supplier, error) {
	items, err := h.svc.List(h.ctx(), search)
	return items, bindError(err)
}

func (h *SupplierHandler) Get(id uint) (*domain.Supplier, error) {
	s, err := h.svc.Get(h.ctx(), id)
	return s, bindError(err)
}

func (h *SupplierHandler) Create(s domain.Supplier) (*domain.Supplier, error) {
	created, err := h.svc.Create(h.ctx(), &s)
	return created, bindError(err)
}

func (h *SupplierHandler) Update(s domain.Supplier) (*domain.Supplier, error) {
	updated, err := h.svc.Update(h.ctx(), &s)
	return updated, bindError(err)
}

func (h *SupplierHandler) Delete(id uint) error {
	return bindError(h.svc.Delete(h.ctx(), id))
}
