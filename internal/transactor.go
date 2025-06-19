package internal

import (
	"context"
	"log"

	"github.com/itsLeonB/ezutil/config"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type GormTransactor struct {
	DB *gorm.DB
}

func (t *GormTransactor) Begin(ctx context.Context) (context.Context, error) {
	tx := t.DB.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return nil, eris.Wrap(err, config.MsgTransactionError)
	}

	return context.WithValue(ctx, config.ContextKeyGormTx, tx), nil
}

func (t *GormTransactor) Commit(ctx context.Context) error {
	tx, err := GetTxFromContext(ctx)
	if err != nil {
		return err
	}
	if tx != nil {
		err = tx.WithContext(ctx).Commit().Error
		if err != nil {
			return eris.Wrap(err, config.MsgTransactionError)
		}
	}

	return nil
}

func (r *GormTransactor) Rollback(ctx context.Context) {
	tx, err := GetTxFromContext(ctx)
	if err != nil {
		log.Println("rollback error:", err)
		return
	}
	if tx == nil {
		log.Println("no transaction is running")
		return
	}

	err = tx.WithContext(ctx).Rollback().Error
	if err != nil {
		if err.Error() == "sql: transaction has already been committed or rolled back" {
			return
		}

		log.Printf("error: %T", err)
		log.Println("rollback error:", err)
	}
}

func (gt *GormTransactor) WithinTransaction(ctx context.Context, serviceFn func(ctx context.Context) error) error {
	ctx, err := gt.Begin(ctx)
	if err != nil {
		return eris.Wrap(err, "error starting transaction")
	}
	defer gt.Rollback(ctx)

	if err := serviceFn(ctx); err != nil {
		return err
	}

	return gt.Commit(ctx)
}

func GetTxFromContext(ctx context.Context) (*gorm.DB, error) {
	trx := ctx.Value(config.ContextKeyGormTx)
	if trx != nil {
		tx, ok := trx.(*gorm.DB)
		if !ok {
			return nil, eris.New("error getting tx from ctx")
		}

		return tx, nil
	}

	return nil, nil
}
