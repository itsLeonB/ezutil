package ezutil

import (
	"time"

	"github.com/itsLeonB/ezutil/internal"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

// Paginate returns a GORM scope that applies pagination to a query.
// It calculates the appropriate offset based on the page number and limit.
// The page parameter is 1-indexed (minimum value of 1).
func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}

		offset := (page - 1) * limit

		return db.Limit(limit).Offset(offset)
	}
}

// OrderBy returns a GORM scope that orders query results by the specified field.
// It uses internal.IsValidFieldName to validate the field name and prevent SQL injection.
// Set ascending to true for ascending order, false for descending.
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

// WhereBySpec returns a GORM scope that applies a WHERE clause based on the provided struct spec.
// Non-zero fields in spec will be used as AND conditions in the query.
func WhereBySpec[T any](spec T) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(&spec)
	}
}

// PreloadRelations returns a GORM scope that preloads the specified relations.
// It eager loads related data to avoid N+1 query problems.
func PreloadRelations(relations []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, relation := range relations {
			db = db.Preload(relation)
		}

		return db
	}
}

func BetweenTime(col string, start, end time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query, _ := GetTimeRangeClause(col, start, end)
		return db.Where(query)
	}
}
