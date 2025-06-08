package ezutil

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil/internal"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

// region Gin Utils

func GetPathParam[T any](ctx *gin.Context, key string) (T, bool, error) {
	var zero T

	paramValue, exists := ctx.Params.Get(key)
	if !exists {
		return zero, false, nil
	}

	parsedValue, err := Parse[T](paramValue)
	if err != nil {
		return zero, false, eris.Wrapf(err, "failed to parse parameter '%s'", key)
	}

	return parsedValue, true, nil
}

func BindRequest[T any](ctx *gin.Context, bindType binding.Binding) (T, error) {
	var zero T

	if err := ctx.ShouldBindWith(&zero, bindType); err != nil {
		return zero, eris.Wrapf(err, "failed to bind request with type %s", bindType.Name())
	}

	return zero, nil
}

// endregion

// region Gorm Utils

func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}

		offset := (page - 1) * limit

		return db.Limit(limit).Offset(offset)
	}
}

func OrderBy(field string, ascending bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Basic validation to prevent SQL injection
		// Only allow alphanumeric characters, underscores, and dots for table.column
		if !internal.IsValidFieldName(field) {
			_ = db.AddError(eris.Errorf("invalid field name: %s", field))
			return db
		}

		if ascending {
			return db.Order(field + " ASC")
		}

		return db.Order(field + " DESC")
	}
}

func WhereBySpec[T any](spec T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(&spec)
	}
}

func PreloadRelations(relations []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, relation := range relations {
			db = db.Preload(relation)
		}

		return db
	}
}

func WithinTransaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return eris.Wrap(tx.Error, "failed to begin transaction")
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// endregion

// region Slice Utils

func MapSlice[T any, U any](input []T, mapperFunc func(T) U) []U {
	output := make([]U, len(input))

	for i, v := range input {
		output[i] = mapperFunc(v)
	}

	return output
}

// endregion

// region Time Utils

func GetStartOfDay(year int, month int, day int) (time.Time, error) {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	// time.Date normalizes invalid dates, so check if the date changed
	if t.Year() != year || int(t.Month()) != month || t.Day() != day {
		return time.Time{}, eris.Errorf("invalid date: %d-%02d-%02d", year, month, day)
	}

	return t, nil
}

func GetEndOfDay(year int, month int, day int) (time.Time, error) {
	t := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.UTC)
	// time.Date normalizes invalid dates, so check if the date changed
	if t.Year() != year || int(t.Month()) != month || t.Day() != day {
		return time.Time{}, eris.Errorf("invalid date: %d-%02d-%02d", year, month, day)
	}

	return t, nil
}

// endregion

// region HTTP Utils

func ServeGracefully(srv *http.Server, timeout time.Duration) {
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("error server listen and serve: %s", err.Error())
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("error shutting down: %s", err.Error())
	}

	log.Println("server successfully shutdown")
}

// endregion

// region String Utils

func Parse[T any](value string) (T, error) {
	var parsed any
	var err error
	var zero T
	var parsedType string

	switch any(zero).(type) {
	case string:
		return any(value).(T), nil
	case int:
		parsed, err = strconv.Atoi(value)
		parsedType = "int"
	case bool:
		parsed, err = strconv.ParseBool(value)
		parsedType = "bool"
	case uuid.UUID:
		parsed, err = uuid.Parse(value)
		parsedType = "uuid"
	default:
		return zero, fmt.Errorf("unsupported type: %T", zero)
	}

	if err != nil {
		return zero, eris.Wrapf(err, "failed to parse value '%s' as %s", value, parsedType)
	}

	return parsed.(T), nil
}

// endregion
