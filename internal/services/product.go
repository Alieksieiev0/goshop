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
	err := pdbs.db.Model(&models.Product{}).Preload("Categories").Find(&products).Error
	return products, err
}

func (pdbs *ProductDBService) CreateProduct(c *gin.Context, product *models.Product) error {
	if len(product.Categories) == 0 {
		return fmt.Errorf("cannot create product without categories")
	}

	return pdbs.db.Create(product).Error
}

func (pdbs *ProductDBService) DeleteAllProducts(c *gin.Context) error {
	return pdbs.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Product{}).Error
}
