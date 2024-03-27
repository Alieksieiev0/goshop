package database

import (
	"fmt"
	"reflect"

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

func ApplyParams(db *gorm.DB, params ...Param) *gorm.DB {
	for _, param := range params {
		db = param(db)
	}
	return db
}

func AppendFilters[T any](params []Param, filters map[string]string) []Param {
	for k, v := range filters {
		t := reflect.TypeOf(*new(T))
		for i := range t.NumField() {
			f := t.Field(i)
			if f.Tag.Get("json") == k {
				params = append(params, Filter(k, v, f.Type != reflect.TypeOf("")))
			}
		}
	}

	return params
}
