package app

import (
	"context"

	"autoshop/internal/domain"
	"autoshop/internal/service"
)

// CustomerHandler is the Wails binding for the Customers module.
type CustomerHandler struct {
	svc *service.CustomerService
}

// NewCustomerHandler constructs the handler.
func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

func (h *CustomerHandler) ctx() context.Context { return context.Background() }

// List returns customers matching the search text (name/phone).
func (h *CustomerHandler) List(search string) ([]domain.Customer, error) {
	items, err := h.svc.List(h.ctx(), search)
	return items, bindError(err)
}

// Get returns one customer.
func (h *CustomerHandler) Get(id uint) (*domain.Customer, error) {
	c, err := h.svc.Get(h.ctx(), id)
	return c, bindError(err)
}

// Create adds a new customer.
func (h *CustomerHandler) Create(c domain.Customer) (*domain.Customer, error) {
	created, err := h.svc.Create(h.ctx(), &c)
	return created, bindError(err)
}

// Update saves changes to a customer.
func (h *CustomerHandler) Update(c domain.Customer) (*domain.Customer, error) {
	updated, err := h.svc.Update(h.ctx(), &c)
	return updated, bindError(err)
}

// Delete soft-deletes a customer.
func (h *CustomerHandler) Delete(id uint) error {
	return bindError(h.svc.Delete(h.ctx(), id))
}

// PurchaseHistory returns a customer's past invoices.
func (h *CustomerHandler) PurchaseHistory(customerID uint) ([]domain.Invoice, error) {
	items, err := h.svc.PurchaseHistory(h.ctx(), customerID)
	return items, bindError(err)
}
