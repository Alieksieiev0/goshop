package services

import (
	"github.com/Alieksieiev0/goshop/internal/models"
	"gorm.io/gorm"
)

type UserService interface {
	Service[models.User]
}

type UserDBService struct {
	DatabaseService[models.User]
}

func NewUserDBService(db *gorm.DB) UserService {
	return &UserDBService{
		DatabaseService: DatabaseService[models.User]{db},
	}
}
