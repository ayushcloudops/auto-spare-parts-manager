package app

import (
	"context"

	"autoshop/internal/domain"
	"autoshop/internal/service"
)

// ProductHandler is the Wails-bound entry point for the Products module. It is
// deliberately thin: it forwards to the service and translates errors for the
// UI. Business logic lives in the service, not here.
type ProductHandler struct {
	svc *service.ProductService
}

// NewProductHandler constructs the handler.
func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// ctx returns a background context. This is a single-user desktop app, so we do
// not thread request-scoped cancellation through the UI boundary.
func (h *ProductHandler) ctx() context.Context { return context.Background() }

// List returns products matching the filter (search/category/low-stock).
func (h *ProductHandler) List(filter domain.ProductFilter) ([]domain.Product, error) {
	items, err := h.svc.List(h.ctx(), filter)
	return items, bindError(err)
}

// Get returns a single product by id.
func (h *ProductHandler) Get(id uint) (*domain.Product, error) {
	p, err := h.svc.Get(h.ctx(), id)
	return p, bindError(err)
}

// Create adds a new product and returns it with its assigned id.
func (h *ProductHandler) Create(p domain.Product) (*domain.Product, error) {
	created, err := h.svc.Create(h.ctx(), &p)
	return created, bindError(err)
}

// Update saves changes to an existing product.
func (h *ProductHandler) Update(p domain.Product) (*domain.Product, error) {
	updated, err := h.svc.Update(h.ctx(), &p)
	return updated, bindError(err)
}

// Delete soft-deletes a product.
func (h *ProductHandler) Delete(id uint) error {
	return bindError(h.svc.Delete(h.ctx(), id))
}

// Categories returns distinct categories for the filter dropdown.
func (h *ProductHandler) Categories() ([]string, error) {
	cats, err := h.svc.Categories(h.ctx())
	return cats, bindError(err)
}

// StockHistory returns the stock ledger for a product.
func (h *ProductHandler) StockHistory(productID uint) ([]domain.StockMovement, error) {
	items, err := h.svc.StockHistory(h.ctx(), productID)
	return items, bindError(err)
}
