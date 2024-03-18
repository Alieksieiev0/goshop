package services

import (
	"github.com/Alieksieiev0/goshop/internal/models"
	"gorm.io/gorm"
)

type CategoryService interface {
	Service[models.Category]
}

type CategoryDBService struct {
	DatabaseService[models.Category]
}

func NewCategoryDBService(db *gorm.DB) CategoryService {
	return &CategoryDBService{
		DatabaseService: DatabaseService[models.Category]{db},
	}
}
