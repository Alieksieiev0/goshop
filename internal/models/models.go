package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Base struct {
	ID        string `gorm:"type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}

type Product struct {
	Base
	Name        string
	Description *string
	Code        string
	Price       decimal.Decimal `sql:"type:decimal(12, 2)"`
	Categories  []*Category     `gorm:"many2many:product_categories;"`
}

type Category struct {
	Base
	Name        string
	Description *string
	ParentId    *uuid.UUID `gorm:"type:uuid;"`
	Parent      *Category
	Products    []*Product `gorm:"many2many:product_categories;"`
}
