package ezutil

import (
	"context"

	"github.com/itsLeonB/ezutil/internal"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

// Transactor provides an interface for managing database transactions with context.
// It abstracts transaction operations to allow for easier testing and different implementations.
type Transactor interface {
	// Begin starts a new database transaction and returns a context containing the transaction.
	Begin(ctx context.Context) (context.Context, error)
	// Commit commits the current transaction in the context.
	Commit(ctx context.Context) error
	// Rollback rolls back the current transaction in the context without returning an error.
	Rollback(ctx context.Context)
}

// NewTransactor creates a new Transactor implementation using GORM.
// The returned Transactor can be used to manage database transactions with context propagation.
func NewTransactor(db *gorm.DB) Transactor {
	return &internal.GormTransactor{DB: db}
}

// GetTxFromContext retrieves the current GORM transaction from the context.
// Returns an error if no transaction is found or if the stored value is not a *gorm.DB.
func GetTxFromContext(ctx context.Context) (*gorm.DB, error) {
	return internal.GetTxFromContext(ctx)
}

// WithinTransaction executes a service function within a database transaction.
// It begins a transaction, executes serviceFn with the transactional context,
// and commits if successful or rolls back if an error occurs.
// serviceFn should use GetTxFromContext to access the transaction.
func WithinTransaction(ctx context.Context, transactor Transactor, serviceFn func(ctx context.Context) error) error {
	ctx, err := transactor.Begin(ctx)
	if err != nil {
		return eris.Wrap(err, "error starting transaction")
	}
	defer transactor.Rollback(ctx)

	if err := serviceFn(ctx); err != nil {
		return eris.Wrap(err, "error executing service function")
	}

	return transactor.Commit(ctx)
}