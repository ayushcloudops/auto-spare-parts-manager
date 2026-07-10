// Package repository contains the GORM-backed implementations of the domain
// repository interfaces, plus shared infrastructure: a generic CRUD base
// (Base[T]) and a context-based transaction mechanism.
package repository

import (
	"context"

	"gorm.io/gorm"
)

// txCtxKey is the unexported key under which an in-flight transaction handle is
// stored in a context. Unexported so nothing outside this package can fabricate
// or read it.
type txCtxKey struct{}

// withTx returns a child context carrying tx. Used by TxManager so that every
// repository call made during a service transaction shares the same *gorm.DB.
func withTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txCtxKey{}, tx)
}

// conn resolves the correct connection for a call: the active transaction if one
// is bound to ctx, otherwise the root db. Either way the request context is
// attached so GORM honours cancellation. This is what makes the unit-of-work
// transparent — repositories never branch on "am I in a transaction?".
func conn(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txCtxKey{}).(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return db.WithContext(ctx)
}
