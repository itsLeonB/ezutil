package ezutil

import (
	"context"

	"github.com/itsLeonB/ezutil/internal"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type Transactor interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context)
}

func NewTransactor(db *gorm.DB) Transactor {
	return &internal.GormTransactor{DB: db}
}

func GetTxFromContext(ctx context.Context) (*gorm.DB, error) {
	return internal.GetTxFromContext(ctx)
}

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
