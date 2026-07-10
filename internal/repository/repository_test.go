package repository_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"autoshop/internal/database"
	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"
	"autoshop/internal/pkg/money"
	"autoshop/internal/repository"

	"gorm.io/gorm"
)

func newDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func sampleProduct() *domain.Product {
	return &domain.Product{
		Name:         "Oil Filter",
		PartNumber:   "OF-2002",
		SellingPrice: money.FromRupees(250),
		GSTRate:      18,
		CurrentStock: 5,
		MinimumStock: 2,
	}
}

func TestBaseCRUD(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewBase[domain.Product](newDB(t))

	p := sampleProduct()
	if err := repo.Create(ctx, p); err != nil {
		t.Fatalf("create: %v", err)
	}
	if p.ID == 0 {
		t.Fatal("expected ID to be set after create")
	}

	got, err := repo.FindByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.Name != "Oil Filter" || got.SellingPrice != money.FromRupees(250) {
		t.Fatalf("unexpected loaded product: %+v", got)
	}

	if err := repo.Delete(ctx, p.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	// After soft-delete it must read as NotFound.
	_, err = repo.FindByID(ctx, p.ID)
	if apperr.KindOf(err) != apperr.KindNotFound {
		t.Fatalf("expected NotFound after delete, got %v", err)
	}
}

func TestFindByIDNotFound(t *testing.T) {
	repo := repository.NewBase[domain.Product](newDB(t))
	_, err := repo.FindByID(context.Background(), 9999)
	if apperr.KindOf(err) != apperr.KindNotFound {
		t.Fatalf("expected NotFound, got %v", err)
	}
}

// TestTransactionRollback proves the unit-of-work: a repo write made inside a
// failing TxManager.Do must not persist.
func TestTransactionRollback(t *testing.T) {
	ctx := context.Background()
	db := newDB(t)
	repo := repository.NewBase[domain.Product](db)
	tx := repository.NewTxManager(db)

	sentinel := errors.New("boom")
	err := tx.Do(ctx, func(ctx context.Context) error {
		if err := repo.Create(ctx, sampleProduct()); err != nil {
			return err
		}
		return sentinel // force rollback
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}

	var count int64
	db.Model(&domain.Product{}).Count(&count)
	if count != 0 {
		t.Fatalf("rollback failed: %d products persisted", count)
	}
}

// TestTransactionCommit proves the happy path commits.
func TestTransactionCommit(t *testing.T) {
	ctx := context.Background()
	db := newDB(t)
	repo := repository.NewBase[domain.Product](db)
	tx := repository.NewTxManager(db)

	err := tx.Do(ctx, func(ctx context.Context) error {
		return repo.Create(ctx, sampleProduct())
	})
	if err != nil {
		t.Fatalf("commit tx: %v", err)
	}

	var count int64
	db.Model(&domain.Product{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 product committed, got %d", count)
	}
}
