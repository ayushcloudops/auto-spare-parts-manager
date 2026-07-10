package database

import (
	"fmt"

	"autoshop/internal/database/migrations"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Migrate applies all pending schema migrations. It is safe to call on every
// startup; gormigrate records applied migration IDs and skips them.
func Migrate(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.All())
	if err := m.Migrate(); err != nil {
		return fmt.Errorf("database: migrate: %w", err)
	}
	return nil
}
