package app

import (
	"context"

	"autoshop/internal/domain"

	"gorm.io/gorm"
)

// Version is the application version surfaced to the UI.
const Version = "0.1.0"

// SystemInfo is returned to the frontend as a health/identity check.
type SystemInfo struct {
	Status   string `json:"status"`   // "ok" | "db_error"
	ShopName string `json:"shopName"` // from the seeded shop profile
	Version  string `json:"version"`
}

// SystemHandler is a Wails-bound handler exposing app-level status. It is the
// reference example of the binding pattern every module follows: a small struct
// holding its dependencies, with exported methods the frontend calls directly.
type SystemHandler struct {
	db *gorm.DB
}

// NewSystemHandler constructs the handler.
func NewSystemHandler(db *gorm.DB) *SystemHandler {
	return &SystemHandler{db: db}
}

// Health reports whether the database is reachable and returns the shop name,
// proving the full frontend → binding → database round-trip works.
func (h *SystemHandler) Health() SystemInfo {
	info := SystemInfo{Status: "ok", Version: Version}

	var profile domain.ShopProfile
	if err := h.db.First(&profile, 1).Error; err == nil {
		info.ShopName = profile.ShopName
	}

	if sqlDB, err := h.db.DB(); err != nil {
		info.Status = "db_error"
	} else if err := sqlDB.PingContext(context.Background()); err != nil {
		info.Status = "db_error"
	}
	return info
}
