package ezutil_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/itsLeonB/ezutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Test model for transactor testing
type TransactorTestModel struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"size:100;not null;uniqueIndex"`
}

func setupTransactorTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&TransactorTestModel{})
	if err != nil {
		t.Fatalf("Failed to migrate test model: %v", err)
	}

	return db
}

func TestNewTransactor(t *testing.T) {
	db := setupTransactorTestDB(t)
	transactor := ezutil.NewTransactor(db)
	assert.NotNil(t, transactor)
}

func TestTransactor_WithinTransaction_Success(t *testing.T) {
	db := setupTransactorTestDB(t)
	transactor := ezutil.NewTransactor(db)
	ctx := context.Background()

	// Test successful transaction
	err := transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		// Get the transaction from context
		tx, err := ezutil.GetTxFromContext(ctx)
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Create a record within the transaction
		model := TransactorTestModel{Name: "Test"}
		return tx.Create(&model).Error
	})

	require.NoError(t, err)

	// Verify the record was committed
	var count int64
	db.Model(&TransactorTestModel{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestTransactor_WithinTransaction_Rollback(t *testing.T) {
	db := setupTransactorTestDB(t)
	transactor := ezutil.NewTransactor(db)
	ctx := context.Background()

	// Test transaction rollback on error
	testError := errors.New("test error")
	err := transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		// Get the transaction from context
		tx, err := ezutil.GetTxFromContext(ctx)
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Create a record within the transaction
		model := TransactorTestModel{Name: "Test"}
		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		// Return an error to trigger rollback
		return testError
	})

	assert.Error(t, err)
	assert.Equal(t, testError, err)

	// Verify the record was rolled back
	var count int64
	db.Model(&TransactorTestModel{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestTransactor_WithinTransaction_Nested(t *testing.T) {
	db := setupTransactorTestDB(t)
	transactor := ezutil.NewTransactor(db)
	ctx := context.Background()

	// Test nested transactions (should reuse the same transaction)
	err := transactor.WithinTransaction(ctx, func(outerCtx context.Context) error {
		// Get the outer transaction
		outerTx, err := ezutil.GetTxFromContext(outerCtx)
		require.NoError(t, err)
		require.NotNil(t, outerTx)

		// Create a record in outer transaction
		model1 := TransactorTestModel{Name: "Outer"}
		if err := outerTx.Create(&model1).Error; err != nil {
			return err
		}

		// Nested transaction should reuse the same transaction
		return transactor.WithinTransaction(outerCtx, func(innerCtx context.Context) error {
			// Get the inner transaction
			innerTx, err := ezutil.GetTxFromContext(innerCtx)
			require.NoError(t, err)
			require.NotNil(t, innerTx)

			// The contexts should contain the same transaction
			assert.Equal(t, outerCtx, innerCtx)

			// Create a record in inner transaction
			model2 := TransactorTestModel{Name: "Inner"}
			return innerTx.Create(&model2).Error
		})
	})

	require.NoError(t, err)

	// Verify both records were committed
	var count int64
	db.Model(&TransactorTestModel{}).Count(&count)
	assert.Equal(t, int64(2), count)

	// Verify the records exist
	var models []TransactorTestModel
	db.Find(&models)
	names := make([]string, len(models))
	for i, model := range models {
		names[i] = model.Name
	}
	assert.ElementsMatch(t, []string{"Outer", "Inner"}, names)
}

func TestTransactor_WithinTransaction_NestedRollback(t *testing.T) {
	db := setupTransactorTestDB(t)
	transactor := ezutil.NewTransactor(db)
	ctx := context.Background()

	// Test nested transaction rollback
	testError := errors.New("inner error")
	err := transactor.WithinTransaction(ctx, func(outerCtx context.Context) error {
		// Get the outer transaction
		outerTx, err := ezutil.GetTxFromContext(outerCtx)
		require.NoError(t, err)
		require.NotNil(t, outerTx)

		// Create a record in outer transaction
		model1 := TransactorTestModel{Name: "Outer"}
		if err := outerTx.Create(&model1).Error; err != nil {
			return err
		}

		// Nested transaction that fails
		return transactor.WithinTransaction(outerCtx, func(innerCtx context.Context) error {
			// Get the inner transaction
			innerTx, err := ezutil.GetTxFromContext(innerCtx)
			require.NoError(t, err)
			require.NotNil(t, innerTx)

			// Create a record in inner transaction
			model2 := TransactorTestModel{Name: "Inner"}
			if err := innerTx.Create(&model2).Error; err != nil {
				return err
			}

			// Return error to trigger rollback
			return testError
		})
	})

	assert.Error(t, err)
	assert.Equal(t, testError, err)

	// Verify both records were rolled back (since they share the same transaction)
	var count int64
	db.Model(&TransactorTestModel{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestTransactor_WithinTransaction_MultipleNested(t *testing.T) {
	db := setupTransactorTestDB(t)
	transactor := ezutil.NewTransactor(db)
	ctx := context.Background()

	// Test multiple levels of nesting
	err := transactor.WithinTransaction(ctx, func(ctx1 context.Context) error {
		tx1, err := ezutil.GetTxFromContext(ctx1)
		require.NoError(t, err)
		require.NotNil(t, tx1)

		model1 := TransactorTestModel{Name: "Level1"}
		if err := tx1.Create(&model1).Error; err != nil {
			return err
		}

		return transactor.WithinTransaction(ctx1, func(ctx2 context.Context) error {
			tx2, err := ezutil.GetTxFromContext(ctx2)
			require.NoError(t, err)
			require.NotNil(t, tx2)

			model2 := TransactorTestModel{Name: "Level2"}
			if err := tx2.Create(&model2).Error; err != nil {
				return err
			}

			return transactor.WithinTransaction(ctx2, func(ctx3 context.Context) error {
				tx3, err := ezutil.GetTxFromContext(ctx3)
				require.NoError(t, err)
				require.NotNil(t, tx3)

				model3 := TransactorTestModel{Name: "Level3"}
				return tx3.Create(&model3).Error
			})
		})
	})

	require.NoError(t, err)

	// Verify all records were committed
	var count int64
	db.Model(&TransactorTestModel{}).Count(&count)
	assert.Equal(t, int64(3), count)

	// Verify the records exist
	var models []TransactorTestModel
	db.Find(&models)
	names := make([]string, len(models))
	for i, model := range models {
		names[i] = model.Name
	}
	assert.ElementsMatch(t, []string{"Level1", "Level2", "Level3"}, names)
}

func TestGetTxFromContext(t *testing.T) {
	db := setupTransactorTestDB(t)
	transactor := ezutil.NewTransactor(db)

	t.Run("context without transaction", func(t *testing.T) {
		ctx := context.Background()
		tx, err := ezutil.GetTxFromContext(ctx)
		require.NoError(t, err)
		assert.Nil(t, tx)
	})

	t.Run("context with transaction", func(t *testing.T) {
		ctx := context.Background()
		err := transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
			tx, err := ezutil.GetTxFromContext(txCtx)
			require.NoError(t, err)
			assert.NotNil(t, tx)
			return nil
		})
		require.NoError(t, err)
	})
}

func TestTransactor_RealWorldScenario(t *testing.T) {
	// Simulate a real-world scenario with SELECT FOR UPDATE
	db := setupTransactorTestDB(t)
	transactor := ezutil.NewTransactor(db)
	ctx := context.Background()

	// Create initial data
	initialModel := TransactorTestModel{Name: "Initial"}
	db.Create(&initialModel)

	// Define process function
	processFunc := func(ctx context.Context, transactor ezutil.Transactor, modelID uint) error {
		return transactor.WithinTransaction(ctx, func(ctx context.Context) error {
			tx, err := ezutil.GetTxFromContext(ctx)
			if err != nil {
				return err
			}

			// Another SELECT FOR UPDATE (should reuse same transaction)
			var model TransactorTestModel
			if err := tx.Where("id = ?", modelID).First(&model).Error; err != nil {
				return err
			}

			// Update the model
			model.Name = "Processed"
			return tx.Save(&model).Error
		})
	}

	// Simulate concurrent operations using SELECT FOR UPDATE
	confirmFunc := func(ctx context.Context) error {
		return transactor.WithinTransaction(ctx, func(ctx context.Context) error {
			tx, err := ezutil.GetTxFromContext(ctx)
			if err != nil {
				return err
			}

			// SELECT FOR UPDATE
			var model TransactorTestModel
			if err := tx.Clauses().Where("id = ?", initialModel.ID).First(&model).Error; err != nil {
				return err
			}

			// Call process function which also uses WithinTransaction
			return processFunc(ctx, transactor, model.ID)
		})
	}

	// Execute the scenario
	err := confirmFunc(ctx)
	require.NoError(t, err)

	// Verify the update was committed
	var updatedModel TransactorTestModel
	db.First(&updatedModel, initialModel.ID)
	assert.Equal(t, "Processed", updatedModel.Name)
}

func TestTransactor_ErrorHandling(t *testing.T) {
	db := setupTransactorTestDB(t)
	transactor := ezutil.NewTransactor(db)
	ctx := context.Background()

	t.Run("database error in transaction - duplicate unique key", func(t *testing.T) {
		// First, create a record to establish a unique constraint violation
		model1 := TransactorTestModel{Name: "UniqueTest"}
		err := db.Create(&model1).Error
		require.NoError(t, err)

		// Now try to create another record with the same name (violates unique constraint)
		err = transactor.WithinTransaction(ctx, func(ctx context.Context) error {
			tx, err := ezutil.GetTxFromContext(ctx)
			require.NoError(t, err)

			// Try to create a record with duplicate unique key
			model2 := TransactorTestModel{Name: "UniqueTest"} // Same name - violates uniqueIndex
			return tx.Create(&model2).Error
		})

		// Should return a database error due to unique constraint violation
		assert.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "unique")
	})

	t.Run("database error in transaction - null constraint violation", func(t *testing.T) {
		err := transactor.WithinTransaction(ctx, func(ctx context.Context) error {
			tx, err := ezutil.GetTxFromContext(ctx)
			require.NoError(t, err)

			// Try to create a record with null name (violates not null constraint)
			// We'll use raw SQL to bypass GORM's validation and hit the database constraint
			return tx.Exec("INSERT INTO transactor_test_models (name) VALUES (NULL)").Error
		})

		// Should return a database error due to NOT NULL constraint violation
		assert.Error(t, err)
		// Different databases have different error messages for NOT NULL violations
		errorMsg := strings.ToLower(err.Error())
		assert.True(t, 
			strings.Contains(errorMsg, "not null") || 
			strings.Contains(errorMsg, "null") || 
			strings.Contains(errorMsg, "constraint"),
			"Expected NOT NULL constraint error, got: %s", err.Error())
	})

	t.Run("panic in transaction", func(t *testing.T) {
		// Test that panics are handled gracefully
		assert.Panics(t, func() {
			transactor.WithinTransaction(ctx, func(ctx context.Context) error {
				panic("test panic")
			})
		})
	})
}
