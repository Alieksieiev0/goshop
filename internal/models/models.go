package models

import (
	"database/sql"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string
	Description sql.NullString
	Code        string
	Price       decimal.Decimal `sql:"type:decimal(12, 2)"`
	Categories  []Category      `gorm:"many2many:product_categories;"`
}

type Category struct {
	gorm.Model
	Name        string
	Description sql.NullString
	ParentId    *uint
	Parent      *Category
	Products    []Product `gorm:"many2many:product_categories;"`
}
