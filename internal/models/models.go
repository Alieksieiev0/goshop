package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Base struct {
	ID        string         `gorm:"type:uuid" json:"id"`
	CreatedAt time.Time      `                 json:"created_at"`
	UpdatedAt time.Time      `                 json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"     json:"deleted_at"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}

type Product struct {
	Base
	Name        string          `json:"name"        gorm:"not null;default:null;"`
	Description *string         `json:"description"`
	Code        string          `json:"code"        gorm:"not null;default:null; unique"`
	Price       decimal.Decimal `json:"price"                                            sql:"type:decimal(12, 2)"`
	Categories  []Category      `json:"categories"  gorm:"many2many:product_categories;"`
}

type Category struct {
	Base
	Name        string     `json:"name"        gorm:"not null; default:null; unique"`
	Description *string    `json:"description"`
	ParentId    *uuid.UUID `json:"parent_id"   gorm:"type:uuid;"`
	Parent      *Category  `json:"parent"`
	Products    []Product  `json:"products"    gorm:"many2many:product_categories;"`
}

type UserRole string

func (ur UserRole) String() string {
	switch ur {
	case Admin:
		return "admin"
	case Usr:
		return "user"
	}
	return "unknown"
}

const (
	Admin UserRole = "admin"
	Usr   UserRole = "user"
)

type User struct {
	Base
	Username string `json:"username" gorm:"default:null;not null;unique;"`
	Email    string `json:"email"    gorm:"default:null;not null;unique;"`
	Password string `json:"password" gorm:"default:null;not null;"`
	Role     UserRole
}
