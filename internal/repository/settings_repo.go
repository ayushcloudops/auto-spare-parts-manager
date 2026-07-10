package repository

import (
	"context"
	"errors"

	"autoshop/internal/domain"
	"autoshop/internal/pkg/apperr"

	"gorm.io/gorm"
)

type settingsRepo struct {
	db *gorm.DB
}

// NewSettingsRepo constructs a SettingsRepository. It manages two entities
// (ShopProfile and AppSetting), so it holds the db directly rather than
// embedding a single-entity Base.
func NewSettingsRepo(db *gorm.DB) domain.SettingsRepository {
	return &settingsRepo{db: db}
}

// GetShopProfile returns the single (ID = 1) shop profile row.
func (r *settingsRepo) GetShopProfile(ctx context.Context) (*domain.ShopProfile, error) {
	var p domain.ShopProfile
	err := conn(ctx, r.db).First(&p, 1).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.NotFound("shop profile not configured")
	}
	if err != nil {
		return nil, apperr.Internal(err, "could not load shop profile")
	}
	return &p, nil
}

// SaveShopProfile upserts the shop profile, always as ID = 1.
func (r *settingsRepo) SaveShopProfile(ctx context.Context, p *domain.ShopProfile) error {
	p.ID = 1
	if err := conn(ctx, r.db).Save(p).Error; err != nil {
		return apperr.Internal(err, "could not save shop profile")
	}
	return nil
}

// Get returns a key/value setting, or "" if absent.
func (r *settingsRepo) Get(ctx context.Context, key string) (string, error) {
	var s domain.AppSetting
	err := conn(ctx, r.db).First(&s, "key = ?", key).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", apperr.Internal(err, "could not read setting")
	}
	return s.Value, nil
}

// Set upserts a key/value setting.
func (r *settingsRepo) Set(ctx context.Context, key, value string) error {
	s := domain.AppSetting{Key: key, Value: value}
	err := conn(ctx, r.db).
		Where(domain.AppSetting{Key: key}).
		Assign(domain.AppSetting{Value: value}).
		FirstOrCreate(&s).Error
	if err != nil {
		return apperr.Internal(err, "could not write setting")
	}
	return nil
}
