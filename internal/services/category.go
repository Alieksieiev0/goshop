package services

import (
	"github.com/Alieksieiev0/goshop/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, category *models.Category) error
}

type CategoryDBService struct {
	db *gorm.DB
}

func NewCategoryDBService(db *gorm.DB) *CategoryDBService {
	return &CategoryDBService{
		db: db,
	}
}

func (cdbs *CategoryDBService) CreateCategory(c *gin.Context, category *models.Category) error {
	result := cdbs.db.Create(category)
	return result.Error
}
