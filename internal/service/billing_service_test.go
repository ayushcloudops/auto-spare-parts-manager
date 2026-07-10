package service_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"autoshop/internal/database"
	"autoshop/internal/domain"
	"autoshop/internal/pkg/money"
	"autoshop/internal/repository"
	"autoshop/internal/service"

	"gorm.io/gorm"
)

// invoiceFY mirrors the repository's financial-year code for the current date,
// so number assertions work regardless of when the tests run.
func invoiceFY() string {
	t := time.Now()
	y := t.Year()
	if int(t.Month()) < int(time.April) {
		y--
	}
	return fmt.Sprintf("%02d%02d", y%100, (y+1)%100)
}

// billingEnv builds a real (temp-file SQLite) billing service with all its
// dependencies, plus helper repos for assertions.
type billingEnv struct {
	db       *gorm.DB
	billing  *service.BillingService
	products domain.ProductRepository
	stock    domain.StockRepository
}

func newBillingEnv(t *testing.T) *billingEnv {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "bill.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := database.Seed(db); err != nil { // seeds shop profile (prefix INV)
		t.Fatalf("seed: %v", err)
	}

	tx := repository.NewTxManager(db)
	products := repository.NewProductRepo(db)
	stock := repository.NewStockRepo(db)
	invoices := repository.NewInvoiceRepo(db)
	customers := repository.NewCustomerRepo(db)
	settings := repository.NewSettingsRepo(db)

	return &billingEnv{
		db:       db,
		billing:  service.NewBillingService(invoices, products, customers, stock, settings, tx),
		products: products,
		stock:    stock,
	}
}

func (e *billingEnv) addProduct(t *testing.T, name string, priceRupees int64, gst float64, stock int) *domain.Product {
	t.Helper()
	p := &domain.Product{Name: name, SellingPrice: money.FromRupees(priceRupees), GSTRate: gst, CurrentStock: stock, MinimumStock: 1}
	if err := e.products.Create(context.Background(), p); err != nil {
		t.Fatalf("add product: %v", err)
	}
	return p
}

func TestCreateBillHappyPath(t *testing.T) {
	env := newBillingEnv(t)
	ctx := context.Background()
	p := env.addProduct(t, "Brake Pad", 1000, 18, 10)

	inv, err := env.billing.Create(ctx, service.CreateBillInput{
		Lines:       []service.BillLineInput{{ProductID: p.ID, Quantity: 2}},
		PaymentMode: domain.PaymentCash,
	})
	if err != nil {
		t.Fatalf("create bill: %v", err)
	}

	// 2 × ₹1000 = ₹2000 taxable, 18% = ₹360 (CGST 180 + SGST 180), grand ₹2360.
	if inv.GrandTotal != money.FromRupees(2360) {
		t.Fatalf("grand total: got %s want 2360.00", inv.GrandTotal)
	}
	if inv.CGST != money.FromRupees(180) || inv.SGST != money.FromRupees(180) {
		t.Fatalf("gst split wrong: %s / %s", inv.CGST, inv.SGST)
	}
	if inv.Number != "INV-"+invoiceFY()+"-0001" {
		t.Fatalf("unexpected invoice number: %s", inv.Number)
	}
	if len(inv.Items) != 1 || inv.Items[0].ProductName != "Brake Pad" {
		t.Fatalf("items not saved/snapshotted: %+v", inv.Items)
	}

	// Stock decremented 10 -> 8.
	reloaded, _ := env.products.FindByID(ctx, p.ID)
	if reloaded.CurrentStock != 8 {
		t.Fatalf("stock: got %d want 8", reloaded.CurrentStock)
	}
	// Ledger has an opening (none here) + a sale movement of -2.
	moves, _ := env.stock.ListByProduct(ctx, p.ID)
	if len(moves) != 1 || moves[0].Delta != -2 || moves[0].Reason != domain.StockReasonSale {
		t.Fatalf("unexpected stock ledger: %+v", moves)
	}

	// Cash sale fully paid.
	if inv.AmountPaid != inv.GrandTotal || inv.AmountDue != money.Zero {
		t.Fatalf("payment: paid %s due %s", inv.AmountPaid, inv.AmountDue)
	}
}

func TestCreateBillInsufficientStockRollsBack(t *testing.T) {
	env := newBillingEnv(t)
	ctx := context.Background()
	p := env.addProduct(t, "Rare Part", 500, 18, 1)

	_, err := env.billing.Create(ctx, service.CreateBillInput{
		Lines:       []service.BillLineInput{{ProductID: p.ID, Quantity: 5}},
		PaymentMode: domain.PaymentCash,
	})
	if err == nil {
		t.Fatal("expected insufficient stock error")
	}

	// Nothing should have changed: stock intact, no invoice, no ledger.
	reloaded, _ := env.products.FindByID(ctx, p.ID)
	if reloaded.CurrentStock != 1 {
		t.Fatalf("stock changed on failed bill: %d", reloaded.CurrentStock)
	}
	var invCount int64
	env.db.Model(&domain.Invoice{}).Count(&invCount)
	if invCount != 0 {
		t.Fatalf("invoice persisted despite failure: %d", invCount)
	}
	moves, _ := env.stock.ListByProduct(ctx, p.ID)
	if len(moves) != 0 {
		t.Fatalf("ledger written despite failure: %+v", moves)
	}
}

func TestCreateBillCreditUpdatesOutstanding(t *testing.T) {
	env := newBillingEnv(t)
	ctx := context.Background()
	p := env.addProduct(t, "Clutch Plate", 1000, 0, 5) // 0% GST for round numbers

	// Create a customer.
	custRepo := repository.NewCustomerRepo(env.db)
	cust := &domain.Customer{Name: "Ravi Motors", CreditLimit: money.FromRupees(10000)}
	if err := custRepo.Create(ctx, cust); err != nil {
		t.Fatalf("create customer: %v", err)
	}

	inv, err := env.billing.Create(ctx, service.CreateBillInput{
		CustomerID:  &cust.ID,
		Lines:       []service.BillLineInput{{ProductID: p.ID, Quantity: 1}},
		PaymentMode: domain.PaymentCredit,
		AmountPaid:  money.Zero,
	})
	if err != nil {
		t.Fatalf("create credit bill: %v", err)
	}
	if inv.AmountDue != money.FromRupees(1000) {
		t.Fatalf("amount due: got %s want 1000.00", inv.AmountDue)
	}

	reloaded, _ := custRepo.FindByID(ctx, cust.ID)
	if reloaded.Outstanding != money.FromRupees(1000) {
		t.Fatalf("outstanding: got %s want 1000.00", reloaded.Outstanding)
	}
}

func TestInvoiceNumbersIncrement(t *testing.T) {
	env := newBillingEnv(t)
	ctx := context.Background()
	p := env.addProduct(t, "Bulb", 100, 18, 100)

	var numbers []string
	for i := 0; i < 3; i++ {
		inv, err := env.billing.Create(ctx, service.CreateBillInput{
			Lines:       []service.BillLineInput{{ProductID: p.ID, Quantity: 1}},
			PaymentMode: domain.PaymentCash,
		})
		if err != nil {
			t.Fatalf("bill %d: %v", i, err)
		}
		numbers = append(numbers, inv.Number)
	}
	fy := invoiceFY()
	want := []string{"INV-" + fy + "-0001", "INV-" + fy + "-0002", "INV-" + fy + "-0003"}
	for i := range want {
		if numbers[i] != want[i] {
			t.Fatalf("number %d: got %s want %s", i, numbers[i], want[i])
		}
	}
}
