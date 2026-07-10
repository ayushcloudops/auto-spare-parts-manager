package service

import (
	"context"
	"strings"
	"time"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"
)

// validGSTRates are the standard Indian GST slabs a product may use.
var validGSTRates = map[float64]bool{0: true, 5: true, 12: true, 18: true, 28: true}

// ProductService implements the product/inventory business rules.
type ProductService struct {
	products domain.ProductRepository
	stock    domain.StockRepository
	tx       Transactor
}

// NewProductService wires the service with its dependencies.
func NewProductService(products domain.ProductRepository, stock domain.StockRepository, tx Transactor) *ProductService {
	return &ProductService{products: products, stock: stock, tx: tx}
}

// validate enforces field-level rules, returning a Validation error listing the
// first problem found.
func (s *ProductService) validate(p *domain.Product) error {
	p.Name = strings.TrimSpace(p.Name)
	switch {
	case p.Name == "":
		return apperr.Validation("product name is required")
	case p.SellingPrice.IsNegative():
		return apperr.Validation("selling price cannot be negative")
	case p.PurchasePrice.IsNegative():
		return apperr.Validation("purchase price cannot be negative")
	case !validGSTRates[p.GSTRate]:
		return apperr.Validation("GST rate must be one of 0, 5, 12, 18 or 28")
	case p.CurrentStock < 0:
		return apperr.Validation("stock cannot be negative")
	case p.MinimumStock < 0:
		return apperr.Validation("minimum stock cannot be negative")
	}
	return nil
}

// Create validates and inserts a new product. If it starts with stock on hand,
// an opening StockMovement is recorded in the same transaction so the ledger is
// always consistent with the product's stock figure.
func (s *ProductService) Create(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	if err := s.validate(p); err != nil {
		return nil, err
	}
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		if err := s.products.Create(ctx, p); err != nil {
			return err
		}
		if p.CurrentStock != 0 {
			return s.stock.Record(ctx, &domain.StockMovement{
				ProductID:    p.ID,
				Delta:        p.CurrentStock,
				BalanceAfter: p.CurrentStock,
				Reason:       domain.StockReasonOpening,
				OccurredAt:   time.Now(),
				Note:         "Opening stock",
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Update validates and saves changes. If the stock figure was edited directly,
// the difference is logged as an adjustment movement, keeping the audit trail
// complete.
func (s *ProductService) Update(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	if p.ID == 0 {
		return nil, apperr.Validation("product id is required")
	}
	if err := s.validate(p); err != nil {
		return nil, err
	}
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		existing, err := s.products.FindByID(ctx, p.ID)
		if err != nil {
			return err
		}
		delta := p.CurrentStock - existing.CurrentStock
		if err := s.products.Update(ctx, p); err != nil {
			return err
		}
		if delta != 0 {
			return s.stock.Record(ctx, &domain.StockMovement{
				ProductID:    p.ID,
				Delta:        delta,
				BalanceAfter: p.CurrentStock,
				Reason:       domain.StockReasonAdjustment,
				OccurredAt:   time.Now(),
				Note:         "Manual edit",
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Get returns a single product.
func (s *ProductService) Get(ctx context.Context, id uint) (*domain.Product, error) {
	return s.products.FindByID(ctx, id)
}

// Delete soft-deletes a product. Past invoices are unaffected (they hold
// snapshots), so this is safe.
func (s *ProductService) Delete(ctx context.Context, id uint) error {
	return s.products.Delete(ctx, id)
}

// List returns products matching the filter.
func (s *ProductService) List(ctx context.Context, f domain.ProductFilter) ([]domain.Product, error) {
	return s.products.List(ctx, f)
}

// Categories returns the distinct non-empty categories, for filter dropdowns.
func (s *ProductService) Categories(ctx context.Context) ([]string, error) {
	return s.products.Categories(ctx)
}

// StockHistory returns the ledger for one product.
func (s *ProductService) StockHistory(ctx context.Context, productID uint) ([]domain.StockMovement, error) {
	return s.stock.ListByProduct(ctx, productID)
}
