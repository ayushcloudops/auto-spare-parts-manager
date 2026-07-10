// Package database opens the SQLite connection and runs schema migrations.
//
// We use the pure-Go SQLite driver (glebarez/sqlite, backed by modernc.org/sqlite)
// so the application cross-compiles to Windows with no cgo / mingw toolchain —
// a hard requirement for shipping a single offline .exe to shop owners.
package database

import (
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open connects to the SQLite database at path and applies pragmatic defaults.
//
// SQLite is single-writer, so we cap the pool at one open connection and enable
// a busy timeout to serialise writes cleanly. WAL mode improves read/write
// concurrency; foreign_keys enforces our cascade constraints.
func Open(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Warn),
		SkipDefaultTransaction: true, // we manage transactions explicitly in services
		NowFunc:                func() time.Time { return time.Now() },
	})
	if err != nil {
		return nil, fmt.Errorf("database: open %q: %w", path, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("database: access sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetConnMaxLifetime(0)

	for _, pragma := range []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		"PRAGMA busy_timeout = 5000",
		"PRAGMA synchronous = NORMAL",
	} {
		if err := db.Exec(pragma).Error; err != nil {
			return nil, fmt.Errorf("database: %q: %w", pragma, err)
		}
	}
	return db, nil
}
