package app

import (
	"context"

	"autoshop/internal/domain"
	"autoshop/internal/service"
)

// SettingsHandler is the Wails binding for the shop profile and app settings.
type SettingsHandler struct {
	svc *service.SettingsService
}

// NewSettingsHandler constructs the handler.
func NewSettingsHandler(svc *service.SettingsService) *SettingsHandler {
	return &SettingsHandler{svc: svc}
}

func (h *SettingsHandler) ctx() context.Context { return context.Background() }

// GetShopProfile returns the shop profile.
func (h *SettingsHandler) GetShopProfile() (*domain.ShopProfile, error) {
	p, err := h.svc.GetShopProfile(h.ctx())
	return p, bindError(err)
}

// SaveShopProfile persists the shop profile.
func (h *SettingsHandler) SaveShopProfile(p domain.ShopProfile) (*domain.ShopProfile, error) {
	saved, err := h.svc.SaveShopProfile(h.ctx(), &p)
	return saved, bindError(err)
}

// GetSetting reads a key/value setting.
func (h *SettingsHandler) GetSetting(key string) (string, error) {
	v, err := h.svc.GetSetting(h.ctx(), key)
	return v, bindError(err)
}

// SetSetting writes a key/value setting.
func (h *SettingsHandler) SetSetting(key, value string) error {
	return bindError(h.svc.SetSetting(h.ctx(), key, value))
}
