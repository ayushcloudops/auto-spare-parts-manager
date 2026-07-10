package service

import (
	"context"
	"time"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"
	"autoshop/internal/pkg/money"
)

// PurchaseLineInput is one received line on a supplier invoice.
type PurchaseLineInput struct {
	ProductID uint        `json:"productId"`
	Quantity  int         `json:"quantity"`
	CostPrice money.Money `json:"costPrice"` // per unit, paise
}

// CreatePurchaseInput is the payload to record a purchase.
type CreatePurchaseInput struct {
	SupplierID    uint                `json:"supplierId"`
	SupplierInvNo string              `json:"supplierInvNo"`
	Lines         []PurchaseLineInput `json:"lines"`
	Notes         string              `json:"notes"`
}

// PurchaseService records stock purchases. Saving a purchase increases each
// product's stock (with a ledger entry) and refreshes its purchase price — all
// atomically.
type PurchaseService struct {
	purchases domain.PurchaseRepository
	products  domain.ProductRepository
	stock     domain.StockRepository
	tx        Transactor
}

// NewPurchaseService wires the service.
func NewPurchaseService(purchases domain.PurchaseRepository, products domain.ProductRepository, stock domain.StockRepository, tx Transactor) *PurchaseService {
	return &PurchaseService{purchases: purchases, products: products, stock: stock, tx: tx}
}

// Create validates and records a purchase, updating stock and cost prices.
func (s *PurchaseService) Create(ctx context.Context, in CreatePurchaseInput) (*domain.Purchase, error) {
	if in.SupplierID == 0 {
		return nil, apperr.Validation("select a supplier")
	}
	if len(in.Lines) == 0 {
		return nil, apperr.Validation("add at least one item")
	}

	var purchaseID uint
	err := s.tx.Do(ctx, func(ctx context.Context) error {
		items := make([]domain.PurchaseItem, len(in.Lines))
		products := make([]*domain.Product, len(in.Lines))
		var subTotal, gstTotal money.Money

		for i, l := range in.Lines {
			if l.Quantity <= 0 {
				return apperr.Validation("quantity must be at least 1")
			}
			p, err := s.products.FindByID(ctx, l.ProductID)
			if err != nil {
				return err
			}
			products[i] = p
			lineNet := l.CostPrice.MulQty(l.Quantity)
			lineGST := lineNet.Percent(p.GSTRate)
			subTotal = subTotal.Add(lineNet)
			gstTotal = gstTotal.Add(lineGST)
			items[i] = domain.PurchaseItem{
				ProductID:   p.ID,
				ProductName: p.Name,
				Quantity:    l.Quantity,
				CostPrice:   l.CostPrice,
				GSTRate:     p.GSTRate,
				LineTotal:   lineNet.Add(lineGST),
			}
		}

		purchase := &domain.Purchase{
			SupplierID:    in.SupplierID,
			SupplierInvNo: in.SupplierInvNo,
			Date:          time.Now(),
			Items:         items,
			SubTotal:      subTotal,
			GSTTotal:      gstTotal,
			GrandTotal:    subTotal.Add(gstTotal),
			Notes:         in.Notes,
		}
		if err := s.purchases.Create(ctx, purchase); err != nil {
			return err
		}
		purchaseID = purchase.ID

		// Increase stock, refresh cost price, append ledger entries.
		for i, l := range in.Lines {
			p := products[i]
			p.CurrentStock += l.Quantity
			p.PurchasePrice = l.CostPrice // latest cost
			if err := s.products.Update(ctx, p); err != nil {
				return err
			}
			if err := s.stock.Record(ctx, &domain.StockMovement{
				ProductID:    p.ID,
				Delta:        l.Quantity,
				BalanceAfter: p.CurrentStock,
				Reason:       domain.StockReasonPurchase,
				RefType:      "purchase",
				RefID:        purchase.ID,
				OccurredAt:   time.Now(),
				Note:         in.SupplierInvNo,
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.purchases.FindByID(ctx, purchaseID)
}

// Get returns a purchase with items + supplier.
func (s *PurchaseService) Get(ctx context.Context, id uint) (*domain.Purchase, error) {
	return s.purchases.FindByID(ctx, id)
}

// List returns purchases matching the filter.
func (s *PurchaseService) List(ctx context.Context, f domain.PurchaseFilter) ([]domain.Purchase, error) {
	return s.purchases.List(ctx, f)
}
