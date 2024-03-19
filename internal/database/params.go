package database

import (
	"fmt"

	"gorm.io/gorm"
)

type Param func(db *gorm.DB) *gorm.DB

func Limit(limit int) Param {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

func Offset(offset int) Param {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

func Order(column string, order string) Param {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(fmt.Sprintf("%s  %s", column, order))
	}
}

func Filter(column string, value string, strict bool) Param {
	return func(db *gorm.DB) *gorm.DB {
		if strict {
			return db.Where(fmt.Sprintf("%s = ?", column), value)
		}
		return db.Where(
			fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", column),
			fmt.Sprintf("%%%s%%", value),
		)
	}
}
