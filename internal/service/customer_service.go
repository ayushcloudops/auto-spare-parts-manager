package service

import (
	"context"
	"strings"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"
)

// CustomerService implements customer management rules.
type CustomerService struct {
	customers domain.CustomerRepository
	invoices  domain.InvoiceRepository
}

// NewCustomerService wires the service.
func NewCustomerService(customers domain.CustomerRepository, invoices domain.InvoiceRepository) *CustomerService {
	return &CustomerService{customers: customers, invoices: invoices}
}

func (s *CustomerService) validate(c *domain.Customer) error {
	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		return apperr.Validation("customer name is required")
	}
	if c.CreditLimit.IsNegative() {
		return apperr.Validation("credit limit cannot be negative")
	}
	return nil
}

func (s *CustomerService) Create(ctx context.Context, c *domain.Customer) (*domain.Customer, error) {
	if err := s.validate(c); err != nil {
		return nil, err
	}
	if err := s.customers.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CustomerService) Update(ctx context.Context, c *domain.Customer) (*domain.Customer, error) {
	if c.ID == 0 {
		return nil, apperr.Validation("customer id is required")
	}
	if err := s.validate(c); err != nil {
		return nil, err
	}
	if err := s.customers.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CustomerService) Get(ctx context.Context, id uint) (*domain.Customer, error) {
	return s.customers.FindByID(ctx, id)
}

func (s *CustomerService) Delete(ctx context.Context, id uint) error {
	return s.customers.Delete(ctx, id)
}

func (s *CustomerService) List(ctx context.Context, search string) ([]domain.Customer, error) {
	return s.customers.List(ctx, search)
}

// PurchaseHistory returns a customer's invoices.
func (s *CustomerService) PurchaseHistory(ctx context.Context, customerID uint) ([]domain.Invoice, error) {
	return s.invoices.List(ctx, domain.InvoiceFilter{CustomerID: &customerID})
}
