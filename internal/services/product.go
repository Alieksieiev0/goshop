package services

import (
	"fmt"

	"github.com/Alieksieiev0/goshop/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductService interface {
	GetAllProducts(c *gin.Context) ([]models.Product, error)
	CreateProduct(c *gin.Context, product *models.Product) error
	DeleteAllProducts(c *gin.Context) error
}

type ProductDBService struct {
	db *gorm.DB
}

func NewProductDBService(db *gorm.DB) ProductService {
	return &ProductDBService{
		db: db,
	}
}

func (pdbs *ProductDBService) GetAllProducts(c *gin.Context) ([]models.Product, error) {
	products := []models.Product{}
	result := pdbs.db.Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf(
			"received error: %v, while executing statement: %v",
			result.Error,
			result.Statement,
		)
	}

	return products, nil
}

func (pdbs *ProductDBService) CreateProduct(c *gin.Context, product *models.Product) error {
	if len(product.Categories) == 0 {
		return fmt.Errorf("cannot create product without categories")
	}

	result := pdbs.db.Create(product)
	if result.Error != nil {
		return fmt.Errorf(
			"received error: %v, while executing statement: %v",
			result.Error,
			result.Statement,
		)
	}
	return nil
}

func (pdbs *ProductDBService) DeleteAllProducts(c *gin.Context) error {
	result := pdbs.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Product{})
	if result.Error != nil {
		return fmt.Errorf(
			"received error: %v, while executing statement: %v",
			result.Error,
			result.Statement,
		)
	}

	return nil
}
