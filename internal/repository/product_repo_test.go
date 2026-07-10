package repository_test

import (
	"context"
	"path/filepath"
	"testing"

	"autoshop/internal/database"
	"autoshop/internal/domain"
	"autoshop/internal/pkg/money"
	"autoshop/internal/repository"

	"gorm.io/gorm"
)

func repoDB(t *testing.T) *gorm.DB {
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

func seedProducts(t *testing.T, repo domain.ProductRepository) {
	t.Helper()
	items := []domain.Product{
		{Name: "Brake Pad", PartNumber: "BP-1", Brand: "Bosch", Category: "Brakes", GSTRate: 28, SellingPrice: money.FromRupees(500), CurrentStock: 2, MinimumStock: 5},
		{Name: "Oil Filter", PartNumber: "OF-9", Brand: "Mann", Category: "Filters", GSTRate: 18, SellingPrice: money.FromRupees(250), CurrentStock: 20, MinimumStock: 5},
		{Name: "Air Filter", PartNumber: "AF-3", Brand: "Bosch", Category: "Filters", GSTRate: 18, SellingPrice: money.FromRupees(300), CurrentStock: 1, MinimumStock: 4},
	}
	for i := range items {
		if err := repo.Create(context.Background(), &items[i]); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
}

func TestProductSearchAndFilters(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewProductRepo(repoDB(t))
	seedProducts(t, repo)

	// Text search matches name/part/brand.
	got, err := repo.List(ctx, domain.ProductFilter{Search: "filter"})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("search 'filter' expected 2, got %d", len(got))
	}

	// Category filter.
	got, _ = repo.List(ctx, domain.ProductFilter{Category: "Brakes"})
	if len(got) != 1 || got[0].Name != "Brake Pad" {
		t.Fatalf("category filter wrong: %+v", got)
	}

	// Low-stock filter: Brake Pad (2<=5) and Air Filter (1<=4) qualify.
	got, _ = repo.List(ctx, domain.ProductFilter{LowStockOnly: true})
	if len(got) != 2 {
		t.Fatalf("low stock expected 2, got %d", len(got))
	}

	low, err := repo.CountLowStock(ctx)
	if err != nil || low != 2 {
		t.Fatalf("CountLowStock expected 2, got %d (err %v)", low, err)
	}
}

func TestProductCategories(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewProductRepo(repoDB(t))
	seedProducts(t, repo)

	cats, err := repo.Categories(ctx)
	if err != nil {
		t.Fatalf("categories: %v", err)
	}
	// Distinct, sorted: Brakes, Filters.
	if len(cats) != 2 || cats[0] != "Brakes" || cats[1] != "Filters" {
		t.Fatalf("unexpected categories: %v", cats)
	}
}
