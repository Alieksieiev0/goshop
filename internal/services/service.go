package services

import (
	"context"
	"fmt"

	"github.com/Alieksieiev0/goshop/internal/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service[T any] interface {
	Get(ctx context.Context, id string) (*T, error)
	GetWithFilters(ctx context.Context, params ...database.Param) (*T, error)
	GetAll(ctx context.Context, params ...database.Param) ([]T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
}

type DatabaseService[T any] struct {
	db *gorm.DB
}

func (dbs *DatabaseService[T]) Get(ctx context.Context, id string) (*T, error) {
	entity := new(T)
	err := dbs.db.Preload(clause.Associations).First(entity, "id = ?", id).Error
	return entity, err
}

func (dbs *DatabaseService[T]) GetWithFilters(
	ctx context.Context,
	params ...database.Param,
) (*T, error) {
	entity := new(T)
	db := applyParams(dbs.db)
	if db == nil {
		return nil, fmt.Errorf("error building query")
	}
	err := dbs.db.Preload(clause.Associations).First(entity).Error
	return entity, err
}

func (dbs *DatabaseService[T]) GetAll(ctx context.Context, params ...database.Param) ([]T, error) {
	entities := []T{}
	db := applyParams(dbs.db)
	if db == nil {
		return nil, fmt.Errorf("error building query")
	}

	err := db.Find(&entities).Error
	return entities, err
}

func (dbs *DatabaseService[T]) Create(ctx context.Context, entity *T) error {
	return dbs.db.Create(entity).Error
}

func (dbs *DatabaseService[T]) Update(ctx context.Context, entity *T) error {
	return dbs.db.Save(entity).Error
}

func (dbs *DatabaseService[T]) Delete(ctx context.Context, id string) error {
	return dbs.db.Delete(new(T), "id = ?", id).Error
}

func applyParams(db *gorm.DB, params ...database.Param) *gorm.DB {
	for _, param := range params {
		db = param(db)
	}
	return db
}
