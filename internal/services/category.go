package services

import (
	"github.com/Alieksieiev0/goshop/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryService interface {
	CreateCategory(c *gin.Context, category *models.Category) error
}

type CategoryDBService struct {
	db *gorm.DB
}

func NewCategoryDBService(db *gorm.DB) CategoryService {
	return &CategoryDBService{
		db: db,
	}
}

func (cdbs *CategoryDBService) CreateCategory(c *gin.Context, category *models.Category) error {
	return cdbs.db.Create(category).Error
}
