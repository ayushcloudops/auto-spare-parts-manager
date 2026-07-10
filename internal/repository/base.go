package repository

import (
	"context"
	"errors"

	"autoshop/internal/pkg/apperr"

	"gorm.io/gorm"
)

// Base is a generic repository providing the CRUD operations common to every
// entity, so concrete repositories (ProductRepo, CustomerRepo, ...) only add the
// queries unique to them. T is the entity type (e.g. domain.Product).
//
// Concrete repos embed Base[T] and use Conn(ctx) for their custom queries, which
// keeps them transaction-aware for free.
type Base[T any] struct {
	db *gorm.DB
}

// NewBase constructs a Base for entity type T.
func NewBase[T any](db *gorm.DB) Base[T] {
	return Base[T]{db: db}
}

// Conn returns the transaction-aware connection for custom queries in concrete
// repositories.
func (r Base[T]) Conn(ctx context.Context) *gorm.DB {
	return conn(ctx, r.db)
}

// Create inserts a new entity.
func (r Base[T]) Create(ctx context.Context, entity *T) error {
	if err := conn(ctx, r.db).Create(entity).Error; err != nil {
		return apperr.Internal(err, "could not create record")
	}
	return nil
}

// Update persists all fields of an existing entity (full save).
func (r Base[T]) Update(ctx context.Context, entity *T) error {
	if err := conn(ctx, r.db).Save(entity).Error; err != nil {
		return apperr.Internal(err, "could not update record")
	}
	return nil
}

// FindByID loads one entity by primary key, returning a typed NotFound error
// when it does not exist.
func (r Base[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := conn(ctx, r.db).First(&entity, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.NotFound("record #%d not found", id)
	}
	if err != nil {
		return nil, apperr.Internal(err, "could not load record")
	}
	return &entity, nil
}

// Delete soft-deletes an entity by primary key (GORM sets DeletedAt). Returns
// NotFound if no row matched.
func (r Base[T]) Delete(ctx context.Context, id uint) error {
	res := conn(ctx, r.db).Delete(new(T), id)
	if res.Error != nil {
		return apperr.Internal(res.Error, "could not delete record")
	}
	if res.RowsAffected == 0 {
		return apperr.NotFound("record #%d not found", id)
	}
	return nil
}
