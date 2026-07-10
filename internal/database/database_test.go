package database_test

import (
	"path/filepath"
	"testing"

	"autoshop/internal/database"
	"autoshop/internal/domain"
	"autoshop/internal/pkg/money"

	"gorm.io/gorm"
)

// newTestDB opens a fresh on-disk SQLite database in a temp dir, migrated and
// seeded. Using a real file (not :memory:) exercises the exact code path used
// in production, including WAL and pragmas.
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	db, err := database.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	if err := database.Seed(db); err != nil {
		t.Fatalf("Seed: %v", err)
	}
	return db
}

func TestMigrateAndSeedAreIdempotent(t *testing.T) {
	db := newTestDB(t)

	// Re-running migrate + seed must not error or duplicate rows.
	if err := database.Migrate(db); err != nil {
		t.Fatalf("second Migrate: %v", err)
	}
	if err := database.Seed(db); err != nil {
		t.Fatalf("second Seed: %v", err)
	}

	var profiles int64
	db.Model(&domain.ShopProfile{}).Count(&profiles)
	if profiles != 1 {
		t.Fatalf("expected exactly 1 shop profile, got %d", profiles)
	}

	var theme domain.AppSetting
	if err := db.First(&theme, "key = ?", domain.SettingTheme).Error; err != nil {
		t.Fatalf("theme setting missing: %v", err)
	}
	if theme.Value != "light" {
		t.Fatalf("expected default theme 'light', got %q", theme.Value)
	}
}

// TestMoneyRoundTrip verifies money.Money persists and reloads exactly, with no
// floating-point drift — the core reason the type exists.
func TestMoneyRoundTrip(t *testing.T) {
	db := newTestDB(t)

	p := domain.Product{
		Name:          "Brake Pad Set",
		PartNumber:    "BP-1001",
		PurchasePrice: money.FromPaise(123455), // ₹1234.55
		SellingPrice:  money.FromPaise(199999), // ₹1999.99
		GSTRate:       28,
		CurrentStock:  10,
		MinimumStock:  3,
	}
	if err := db.Create(&p).Error; err != nil {
		t.Fatalf("create product: %v", err)
	}

	var got domain.Product
	if err := db.First(&got, p.ID).Error; err != nil {
		t.Fatalf("reload product: %v", err)
	}
	if got.SellingPrice != money.FromPaise(199999) {
		t.Fatalf("selling price drift: got %d paise", got.SellingPrice.Paise())
	}
	if got.PurchasePrice.Paise() != 123455 {
		t.Fatalf("purchase price drift: got %d paise", got.PurchasePrice.Paise())
	}
	if !got.IsLowStock() == (got.CurrentStock <= got.MinimumStock) {
		// sanity: low-stock logic consistent
		t.Fatalf("low stock logic inconsistent")
	}
}
