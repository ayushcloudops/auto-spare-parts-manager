package repository

import (
	"context"

	"gorm.io/gorm"
)

// TxManager runs a function inside a single database transaction.
//
// Services use it to make multi-step operations atomic. The classic example is
// saving a bill: allocate invoice number, write invoice + items, decrement
// stock, append StockMovement rows, update customer outstanding — if any step
// fails, the whole thing rolls back and the database is untouched.
//
// The transaction handle is injected into the context, so every repository
// invoked inside fn automatically enrols in the same transaction with no extra
// plumbing.
type TxManager struct {
	db *gorm.DB
}

// NewTxManager constructs a TxManager bound to the root connection.
func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

// Do executes fn within a transaction. Returning a non-nil error rolls back;
// returning nil commits.
func (m *TxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(withTx(ctx, tx))
	})
}
