package ezutil

import (
	"github.com/itsLeonB/ezutil/internal"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

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
