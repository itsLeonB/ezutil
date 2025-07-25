package ezutil

import (
	"context"
	"reflect"

	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

// CRUDRepository defines a generic interface for basic CRUD operations on entities of type T.
// It provides standard database operations with context support and transaction awareness.
// The interface abstracts the underlying database implementation for easier testing and flexibility.
type CRUDRepository[T any] interface {
	// Insert creates a new record in the database.
	Insert(ctx context.Context, model T) (T, error)
	// FindAll retrieves multiple records based on the specification.
	FindAll(ctx context.Context, spec Specification[T]) ([]T, error)
	// FindFirst retrieves the first record matching the specification.
	FindFirst(ctx context.Context, spec Specification[T]) (T, error)
	// Update modifies an existing record in the database.
	Update(ctx context.Context, model T) (T, error)
	// Delete removes a record from the database (hard delete).
	Delete(ctx context.Context, model T) error
	// BatchInsert creates multiple records in a single database operation.
	BatchInsert(ctx context.Context, models []T) ([]T, error)
	// GetGormInstance returns the appropriate GORM DB instance (transaction-aware).
	GetGormInstance(ctx context.Context) (*gorm.DB, error)
}

// Specification defines query parameters for database operations.
// It includes the model for WHERE conditions, relations to preload, and locking options.
type Specification[T any] struct {
	Model            T        // Model with fields set for WHERE conditions
	PreloadRelations []string // Relations to eager load
	ForUpdate        bool     // Whether to use SELECT ... FOR UPDATE
}

// NewCRUDRepository creates a new CRUD repository implementation using GORM.
// The repository provides transaction-aware database operations for the specified entity type T.
func NewCRUDRepository[T any](db *gorm.DB) CRUDRepository[T] {
	return &crudRepositoryGorm[T]{db}
}

type crudRepositoryGorm[T any] struct {
	db *gorm.DB
}

func (cr *crudRepositoryGorm[T]) Insert(ctx context.Context, model T) (T, error) {
	var zero T

	if err := cr.checkZeroValue(model); err != nil {
		return zero, err
	}

	db, err := cr.GetGormInstance(ctx)
	if err != nil {
		return zero, err
	}

	if err = db.Create(&model).Error; err != nil {
		return zero, eris.Wrap(err, "error inserting data")
	}

	return model, nil
}

func (cr *crudRepositoryGorm[T]) FindAll(ctx context.Context, spec Specification[T]) ([]T, error) {
	var models []T

	db, err := cr.GetGormInstance(ctx)
	if err != nil {
		return nil, err
	}

	err = db.Scopes(
		WhereBySpec(spec.Model),
		DefaultOrder(),
		PreloadRelations(spec.PreloadRelations),
		ForUpdate(spec.ForUpdate),
	).
		Find(&models).
		Error

	if err != nil {
		return nil, eris.Wrap(err, "error querying data")
	}

	return models, nil
}

func (cr *crudRepositoryGorm[T]) FindFirst(ctx context.Context, spec Specification[T]) (T, error) {
	var model T

	db, err := cr.GetGormInstance(ctx)
	if err != nil {
		return model, err
	}

	err = db.Scopes(
		WhereBySpec(spec.Model),
		DefaultOrder(),
		PreloadRelations(spec.PreloadRelations),
		ForUpdate(spec.ForUpdate),
	).
		First(&model).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model, nil
		}
		return model, eris.Wrap(err, "error querying data")
	}

	return model, nil
}

func (cr *crudRepositoryGorm[T]) Update(ctx context.Context, model T) (T, error) {
	var zero T

	if err := cr.checkZeroValue(model); err != nil {
		return zero, err
	}

	db, err := cr.GetGormInstance(ctx)
	if err != nil {
		return zero, err
	}

	if err = db.Save(&model).Error; err != nil {
		return zero, eris.Wrap(err, "error updating data")
	}

	return model, nil
}

func (cr *crudRepositoryGorm[T]) Delete(ctx context.Context, model T) error {
	if err := cr.checkZeroValue(model); err != nil {
		return err
	}

	db, err := cr.GetGormInstance(ctx)
	if err != nil {
		return err
	}

	if err = db.Unscoped().Delete(&model).Error; err != nil {
		return eris.Wrap(err, "error deleting data")
	}

	return nil
}

func (cr *crudRepositoryGorm[T]) BatchInsert(ctx context.Context, models []T) ([]T, error) {
	if len(models) < 1 {
		return nil, eris.Errorf("inserted models cannot be empty")
	}

	db, err := cr.GetGormInstance(ctx)
	if err != nil {
		return nil, err
	}

	if err = db.Create(&models).Error; err != nil {
		return nil, eris.Wrap(err, "error batch inserting data")
	}

	return models, nil
}

func (cr *crudRepositoryGorm[T]) checkZeroValue(model T) error {
	if reflect.DeepEqual(model, *new(T)) {
		return eris.New("model cannot be zero value")
	}

	return nil
}

func (cr *crudRepositoryGorm[T]) GetGormInstance(ctx context.Context) (*gorm.DB, error) {
	tx, err := GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx, nil
	}

	return cr.db.WithContext(ctx), nil
}
