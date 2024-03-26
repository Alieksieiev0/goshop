package repositories

import (
	"context"
	"fmt"

	"github.com/Alieksieiev0/goshop/internal/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository[T any] interface {
	Get(ctx context.Context, id string) (*T, error)
	Save(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
}

type DatabaseRepository[T any] interface {
	Repository[T]
	GetWithFilters(ctx context.Context, params ...database.Param) (*T, error)
	GetAll(ctx context.Context, params ...database.Param) ([]T, error)
}

type GormRepository[T any] struct {
	db *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{
		db: db,
	}
}

func (gr *GormRepository[T]) Get(ctx context.Context, id string) (*T, error) {
	entity := new(T)
	err := gr.db.Preload(clause.Associations).First(entity, "id = ?", id).Error
	return entity, err
}

func (gr *GormRepository[T]) GetWithFilters(
	ctx context.Context,
	params ...database.Param,
) (*T, error) {
	entity := new(T)
	db := database.ApplyParams(gr.db, params...)
	if db == nil {
		return nil, fmt.Errorf("error building query")
	}
	err := db.Preload(clause.Associations).First(entity).Error
	return entity, err
}

func (gr *GormRepository[T]) GetAll(ctx context.Context, params ...database.Param) ([]T, error) {
	entities := []T{}
	err := gr.db.Find(&entities).Error
	return entities, err
}

func (gr *GormRepository[T]) Save(ctx context.Context, entity *T) error {
	return gr.db.Save(entity).Error
}

func (gr *GormRepository[T]) Delete(ctx context.Context, id string) error {
	return gr.db.Delete(new(T), "id = ?", id).Error
}
