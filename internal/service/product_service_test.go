package service_test

import (
	"context"
	"testing"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"
	"autoshop/internal/pkg/money"
	"autoshop/internal/service"
)

// --- fakes -----------------------------------------------------------------
// These in-memory fakes let us unit-test the service with no database, which is
// exactly the payoff of depending on repository interfaces.

type fakeProductRepo struct {
	items  map[uint]*domain.Product
	nextID uint
}

func newFakeProductRepo() *fakeProductRepo {
	return &fakeProductRepo{items: map[uint]*domain.Product{}, nextID: 1}
}

func (f *fakeProductRepo) Create(_ context.Context, p *domain.Product) error {
	p.ID = f.nextID
	f.nextID++
	cp := *p
	f.items[p.ID] = &cp
	return nil
}
func (f *fakeProductRepo) Update(_ context.Context, p *domain.Product) error {
	if _, ok := f.items[p.ID]; !ok {
		return apperr.NotFound("record #%d not found", p.ID)
	}
	cp := *p
	f.items[p.ID] = &cp
	return nil
}
func (f *fakeProductRepo) FindByID(_ context.Context, id uint) (*domain.Product, error) {
	if p, ok := f.items[id]; ok {
		cp := *p
		return &cp, nil
	}
	return nil, apperr.NotFound("record #%d not found", id)
}
func (f *fakeProductRepo) Delete(_ context.Context, id uint) error {
	if _, ok := f.items[id]; !ok {
		return apperr.NotFound("record #%d not found", id)
	}
	delete(f.items, id)
	return nil
}
func (f *fakeProductRepo) List(_ context.Context, _ domain.ProductFilter) ([]domain.Product, error) {
	out := make([]domain.Product, 0, len(f.items))
	for _, p := range f.items {
		out = append(out, *p)
	}
	return out, nil
}
func (f *fakeProductRepo) Count(_ context.Context, _ domain.ProductFilter) (int64, error) {
	return int64(len(f.items)), nil
}
func (f *fakeProductRepo) CountAll(_ context.Context) (int64, error) {
	return int64(len(f.items)), nil
}
func (f *fakeProductRepo) CountLowStock(_ context.Context) (int64, error) { return 0, nil }
func (f *fakeProductRepo) Categories(_ context.Context) ([]string, error) { return nil, nil }

type fakeStockRepo struct{ movements []domain.StockMovement }

func (f *fakeStockRepo) Record(_ context.Context, m *domain.StockMovement) error {
	f.movements = append(f.movements, *m)
	return nil
}
func (f *fakeStockRepo) ListByProduct(_ context.Context, id uint) ([]domain.StockMovement, error) {
	var out []domain.StockMovement
	for _, m := range f.movements {
		if m.ProductID == id {
			out = append(out, m)
		}
	}
	return out, nil
}

// noopTx runs fn directly (no real transaction needed for unit tests).
type noopTx struct{}

func (noopTx) Do(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

// --- tests -----------------------------------------------------------------

func newSvc() (*service.ProductService, *fakeProductRepo, *fakeStockRepo) {
	pr := newFakeProductRepo()
	sr := &fakeStockRepo{}
	return service.NewProductService(pr, sr, noopTx{}), pr, sr
}

func TestCreateValidatesName(t *testing.T) {
	svc, _, _ := newSvc()
	_, err := svc.Create(context.Background(), &domain.Product{Name: "   ", GSTRate: 18})
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestCreateRejectsBadGST(t *testing.T) {
	svc, _, _ := newSvc()
	_, err := svc.Create(context.Background(), &domain.Product{Name: "Filter", GSTRate: 17})
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("expected validation error for GST 17, got %v", err)
	}
}

func TestCreateRecordsOpeningStock(t *testing.T) {
	svc, _, stock := newSvc()
	p, err := svc.Create(context.Background(), &domain.Product{
		Name:         "Spark Plug",
		GSTRate:      28,
		SellingPrice: money.FromRupees(120),
		CurrentStock: 15,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if p.ID == 0 {
		t.Fatal("expected id assigned")
	}
	if len(stock.movements) != 1 {
		t.Fatalf("expected 1 opening movement, got %d", len(stock.movements))
	}
	m := stock.movements[0]
	if m.Delta != 15 || m.Reason != domain.StockReasonOpening {
		t.Fatalf("unexpected opening movement: %+v", m)
	}
}

func TestUpdateLogsStockAdjustment(t *testing.T) {
	svc, _, stock := newSvc()
	p, _ := svc.Create(context.Background(), &domain.Product{Name: "Belt", GSTRate: 18, CurrentStock: 10})
	stock.movements = nil // clear the opening movement

	p.CurrentStock = 7 // sold/adjusted down by 3
	if _, err := svc.Update(context.Background(), p); err != nil {
		t.Fatalf("update: %v", err)
	}
	if len(stock.movements) != 1 {
		t.Fatalf("expected 1 adjustment movement, got %d", len(stock.movements))
	}
	if stock.movements[0].Delta != -3 || stock.movements[0].Reason != domain.StockReasonAdjustment {
		t.Fatalf("unexpected adjustment: %+v", stock.movements[0])
	}
}
