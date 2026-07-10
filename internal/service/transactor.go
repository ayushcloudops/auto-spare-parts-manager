// Package service holds the business logic layer. Services orchestrate
// repositories, enforce validation and invariants, and own transaction
// boundaries. They depend on domain interfaces (repositories) and a Transactor,
// never on GORM directly — which keeps them unit-testable with fakes.
package service

import "context"

// Transactor runs a function inside a single database transaction. The
// repository.TxManager satisfies this; declaring it here as an interface keeps
// the service layer decoupled from the concrete infrastructure.
type Transactor interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
