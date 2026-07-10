package database

import (
	"fmt"

	"autoshop/internal/domain"

	"gorm.io/gorm"
)

// Seed ensures baseline rows exist so the app is usable on first launch. It is
// idempotent: running it repeatedly never duplicates or overwrites data the
// shopkeeper has since edited.
func Seed(db *gorm.DB) error {
	if err := seedShopProfile(db); err != nil {
		return err
	}
	return seedDefaultSettings(db)
}

// seedShopProfile inserts a placeholder shop profile (ID = 1) the owner edits
// in Settings. We only create it if absent.
func seedShopProfile(db *gorm.DB) error {
	var count int64
	if err := db.Model(&domain.ShopProfile{}).Count(&count).Error; err != nil {
		return fmt.Errorf("database: count shop profile: %w", err)
	}
	if count > 0 {
		return nil
	}
	profile := domain.ShopProfile{
		ShopName:      "My Auto Spare Parts",
		InvoicePrefix: "INV",
		ReceiptFooter: "Thank You Visit Again",
	}
	profile.ID = 1
	if err := db.Create(&profile).Error; err != nil {
		return fmt.Errorf("database: seed shop profile: %w", err)
	}
	return nil
}

// seedDefaultSettings inserts default key/value settings if missing.
func seedDefaultSettings(db *gorm.DB) error {
	defaults := []domain.AppSetting{
		{Key: domain.SettingTheme, Value: "light"},
	}
	for _, s := range defaults {
		// FirstOrCreate keeps any existing user-set value untouched.
		if err := db.Where(domain.AppSetting{Key: s.Key}).
			FirstOrCreate(&s).Error; err != nil {
			return fmt.Errorf("database: seed setting %q: %w", s.Key, err)
		}
	}
	return nil
}
