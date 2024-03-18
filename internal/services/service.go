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

func (dbs *DatabaseService[T]) GetAll(ctx context.Context, params ...database.Param) ([]T, error) {
	entities := []T{}
	var db *gorm.DB
	for i, param := range params {
		if i == 0 {
			db = param(dbs.db)
			continue
		}
		db = param(db)
	}

	if db == nil {
		return nil, fmt.Errorf("no params were provided to build query")
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
	return dbs.db.Delete(new(T), id).Error
}
