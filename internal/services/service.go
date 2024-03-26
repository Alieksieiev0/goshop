package services

import (
	"context"

	"github.com/Alieksieiev0/goshop/internal/models"
	"github.com/Alieksieiev0/goshop/internal/repositories"
	"gorm.io/gorm"
)

type Service[T any] interface {
	GetById(ctx context.Context, id string) (*T, error)
	SaveEntity(ctx context.Context, entity *T) error
	DeleteById(ctx context.Context, id string) error
}

type CrudService[T any] struct {
	repositories repositories.Repository[T]
}

func NewCrudService[T any](repositories repositories.Repository[T]) *CrudService[T] {
	return &CrudService[T]{
		repositories: repositories,
	}
}

func (cs *CrudService[T]) GetById(ctx context.Context, id string) (*T, error) {
	return cs.repositories.Get(ctx, id)
}

func (cs *CrudService[T]) SaveEntity(ctx context.Context, entity *T) error {
	return cs.repositories.Save(ctx, entity)
}

func (cs *CrudService[T]) DeleteById(ctx context.Context, id string) error {
	return cs.repositories.Delete(ctx, id)
}

type CategoryService interface {
	Service[models.Category]
}

type CategoryDatabaseService struct {
	*CrudService[models.Category]
}

func NewCategoryDatabaseService(db *gorm.DB) CategoryService {
	return &CategoryDatabaseService{
		CrudService: NewCrudService(repositories.NewGormRepository[models.Category](db)),
	}
}

type ProductService interface {
	Service[models.Product]
}

type ProductDatabaseService struct {
	*CrudService[models.Product]
}

func NewProductDatabaseService(db *gorm.DB) ProductService {
	return &ProductDatabaseService{
		CrudService: NewCrudService(repositories.NewGormRepository[models.Product](db)),
	}
}

type UserService interface {
	Service[models.User]
}

type UserDatabaseService struct {
	*CrudService[models.User]
}

func NewUserDatabaseService(db *gorm.DB) UserService {
	return &UserDatabaseService{
		CrudService: NewCrudService(repositories.NewGormRepository[models.User](db)),
	}
}

type AuthService interface {
	Register()
	Login()
}
