package service

import (
	"context"
	"strings"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"
)

// SupplierService implements supplier management.
type SupplierService struct {
	suppliers domain.SupplierRepository
}

// NewSupplierService wires the service.
func NewSupplierService(suppliers domain.SupplierRepository) *SupplierService {
	return &SupplierService{suppliers: suppliers}
}

func (s *SupplierService) validate(sup *domain.Supplier) error {
	sup.Name = strings.TrimSpace(sup.Name)
	if sup.Name == "" {
		return apperr.Validation("supplier name is required")
	}
	return nil
}

func (s *SupplierService) Create(ctx context.Context, sup *domain.Supplier) (*domain.Supplier, error) {
	if err := s.validate(sup); err != nil {
		return nil, err
	}
	if err := s.suppliers.Create(ctx, sup); err != nil {
		return nil, err
	}
	return sup, nil
}

func (s *SupplierService) Update(ctx context.Context, sup *domain.Supplier) (*domain.Supplier, error) {
	if sup.ID == 0 {
		return nil, apperr.Validation("supplier id is required")
	}
	if err := s.validate(sup); err != nil {
		return nil, err
	}
	if err := s.suppliers.Update(ctx, sup); err != nil {
		return nil, err
	}
	return sup, nil
}

func (s *SupplierService) Get(ctx context.Context, id uint) (*domain.Supplier, error) {
	return s.suppliers.FindByID(ctx, id)
}

func (s *SupplierService) Delete(ctx context.Context, id uint) error {
	return s.suppliers.Delete(ctx, id)
}

func (s *SupplierService) List(ctx context.Context, search string) ([]domain.Supplier, error) {
	return s.suppliers.List(ctx, search)
}
