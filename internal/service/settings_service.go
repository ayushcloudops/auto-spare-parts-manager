package service

import (
	"context"
	"strings"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"
)

// SettingsService manages the shop profile and app settings.
type SettingsService struct {
	settings domain.SettingsRepository
}

// NewSettingsService wires the service.
func NewSettingsService(settings domain.SettingsRepository) *SettingsService {
	return &SettingsService{settings: settings}
}

// GetShopProfile returns the shop profile.
func (s *SettingsService) GetShopProfile(ctx context.Context) (*domain.ShopProfile, error) {
	return s.settings.GetShopProfile(ctx)
}

// SaveShopProfile validates and persists the shop profile.
func (s *SettingsService) SaveShopProfile(ctx context.Context, p *domain.ShopProfile) (*domain.ShopProfile, error) {
	p.ShopName = strings.TrimSpace(p.ShopName)
	if p.ShopName == "" {
		return nil, apperr.Validation("shop name is required")
	}
	if strings.TrimSpace(p.InvoicePrefix) == "" {
		p.InvoicePrefix = "INV"
	}
	if err := s.settings.SaveShopProfile(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// GetSetting reads a key/value setting.
func (s *SettingsService) GetSetting(ctx context.Context, key string) (string, error) {
	return s.settings.Get(ctx, key)
}

// SetSetting writes a key/value setting.
func (s *SettingsService) SetSetting(ctx context.Context, key, value string) error {
	return s.settings.Set(ctx, key, value)
}
