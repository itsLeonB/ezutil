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
	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil/internal"
	"github.com/itsLeonB/ezutil/types"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

// region Types

type QueryOptions struct {
	Page  int
	Limit int
}

const (
	PathParam  = types.ParamTypePath
	QueryParam = types.ParamTypeQuery

	BindJSON  = types.BindTypeJSON
	BindForm  = types.BindTypeForm
	BindQuery = types.BindTypeQuery
)

// endregion

// GetParam retrieves a parameter from the Gin context by type and key, parses it into type T, and returns the value, a boolean indicating existence, and an error if parsing fails.

func GetParam[T any](ctx *gin.Context, paramType types.ParamType, key string) (T, bool, error) {
	var zero T

	paramValue, exists := internal.GetParamByType(ctx, paramType, key)
	if !exists {
		return zero, false, nil
	}

	parsedValue, err := Parse[T](paramValue)
	if err != nil {
		return zero, false, eris.Wrapf(err, "failed to parse parameter '%s'", key)
	}

	return parsedValue, true, nil
}

// GetPagination extracts pagination parameters from the Gin context, applying defaults if missing or invalid.
// Returns a QueryOptions struct with validated page and limit values, or an error if parsing fails.
func GetPagination(ctx *gin.Context, defaultLimit int) (QueryOptions, error) {
	page, exists, err := GetParam[int](ctx, QueryParam, "page")
	if err != nil {
		return QueryOptions{}, eris.Wrapf(err, "failed to get 'page' parameter")
	}
	if !exists || page < 1 {
		page = 1 // Default page
	}

	limit, exists, err := GetParam[int](ctx, QueryParam, "limit")
	if err != nil {
		return QueryOptions{}, eris.Wrapf(err, "failed to get 'limit' parameter")
	}
	if !exists || limit < 1 || defaultLimit < 1 {
		limit = defaultLimit // Default limit
	}

	return QueryOptions{Page: page, Limit: limit}, nil
}

// BindRequest binds request data from the Gin context to a struct of type T using the specified binding type.
// Returns the populated struct or an error if binding fails or the bind type is unsupported.
func BindRequest[T any](ctx *gin.Context, bindType types.BindType) (T, error) {
	var zero T

	switch bindType {
	case types.BindTypeJSON:
		if err := ctx.ShouldBindJSON(&zero); err != nil {
			return zero, err
		}
	case types.BindTypeForm:
		if err := ctx.ShouldBind(&zero); err != nil {
			return zero, err
		}
	default:
		return zero, eris.Errorf("unsupported bind type: %s", bindType)
	}

	return zero, nil
}

// endregion

// Paginate returns a Gorm scope function that applies pagination using the specified page and limit. Page numbers less than 1 are treated as 1.

func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}

		offset := (page - 1) * limit

		return db.Limit(limit).Offset(offset)
	}
}

// OrderBy returns a Gorm scope function that orders query results by the specified field in ascending or descending order.
func OrderBy(field string, ascending bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if ascending {
			return db.Order(field + " ASC")
		}

		return db.Order(field + " DESC")
	}
}

// WhereBySpec returns a Gorm scope function that applies a WHERE clause using the provided spec struct.
func WhereBySpec[T any](db *gorm.DB, spec T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(&spec)
	}
}

// PreloadRelations returns a Gorm scope function that preloads the specified relations for a query.
func PreloadRelations(relations []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, relation := range relations {
			db = db.Preload(relation)
		}

		return db
	}
}

// WithinTransaction executes the provided function within a database transaction, rolling back on error and committing on success. Returns an error if the transaction fails.
func WithinTransaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Begin()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction failed: %w", err)
	}

	tx.Commit()

	return nil
}

// endregion

// MapSlice applies a mapping function to each element of the input slice and returns a new slice with the mapped values.

func MapSlice[T any, U any](input []T, mapperFunc func(T) U) []U {
	output := make([]U, len(input))

	for i, v := range input {
		output[i] = mapperFunc(v)
	}

	return output
}

// endregion

// GetStartOfDay returns the UTC time at 00:00:00 for the specified date.
// Returns an error if the provided year, month, or day is invalid or if the date cannot be parsed.

func GetStartOfDay(year int, month int, day int) (time.Time, error) {
	if year < 1970 || month < 1 || month > 12 || day < 1 || day > 31 {
		return time.Time{}, eris.Errorf("invalid date: %d-%02d-%02d", year, month, day)
	}

	startOfDay := fmt.Sprintf("%04d-%02d-%02dT00:00:00Z", year, month, day)
	t, err := time.Parse(time.RFC3339, startOfDay)
	if err != nil {
		return time.Time{}, eris.Wrapf(err, "failed to parse date: %s", startOfDay)
	}

	return t, nil
}

// GetEndOfDay returns a UTC time representing 23:59:59 on the specified date.
// Returns an error if the provided date components are invalid or cannot be parsed.
func GetEndOfDay(year int, month int, day int) (time.Time, error) {
	if year < 1970 || month < 1 || month > 12 || day < 1 || day > 31 {
		return time.Time{}, eris.Errorf("invalid date: %d-%02d-%02d", year, month, day)
	}

	endOfDay := fmt.Sprintf("%04d-%02d-%02dT23:59:59Z", year, month, day)
	t, err := time.Parse(time.RFC3339, endOfDay)
	if err != nil {
		return time.Time{}, eris.Wrapf(err, "failed to parse date: %s", endOfDay)
	}

	return t, nil
}

// endregion

// ServeGracefully starts the HTTP server and shuts it down gracefully on receiving an interrupt or termination signal, allowing active connections to complete within the specified timeout.

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

// Parse converts a string to the specified type T, supporting string, int, bool, and uuid.UUID.
// Returns an error if parsing fails or if the type is unsupported.

func Parse[T any](value string) (T, error) {
	var zero T

	switch any(zero).(type) {
	case string:
		return any(value).(T), nil
	case int:
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return zero, err
		}

		return any(parsed).(T), nil
	case bool:
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return zero, err
		}

		return any(parsed).(T), nil
	case uuid.UUID:
		parsed, err := uuid.Parse(value)
		if err != nil {
			return zero, err
		}

		return any(parsed).(T), nil
	default:
		return zero, fmt.Errorf("unsupported type: %T", zero)
	}
}

// endregion
