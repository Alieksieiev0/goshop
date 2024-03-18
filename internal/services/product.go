package services

import (
	"context"

	"github.com/Alieksieiev0/goshop/internal/models"
	"gorm.io/gorm"
)

type ProductService interface {
	Service[models.Product]
	DeleteAllProducts(ctx context.Context) error
}

type ProductDBService struct {
	DatabaseService[models.Product]
}

func NewProductDBService(db *gorm.DB) ProductService {
	return &ProductDBService{
		DatabaseService: DatabaseService[models.Product]{db},
	}
}

func (pdbs *ProductDBService) DeleteAllProducts(ctx context.Context) error {
	return pdbs.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Product{}).Error
}
