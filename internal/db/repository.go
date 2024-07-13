package db

import (
	"context"

	"gorm.io/gorm"
)

type IRepository[T any] interface{
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	GetAll(ctx context.Context) ([]T, error)
	GetById(ctx context.Context, id string) (T, error)
	GetByName(ctx context.Context, name string) (T, error)
}

type Repository[T any] struct {
	database *gorm.DB
}

func New[T any](database *gorm.DB) *Repository[T] {
	return &Repository[T]{
		database: database,
	}
}

func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	return r.database.WithContext(ctx).Create(entity).Error
}

func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	return r.database.WithContext(ctx).Updates(entity).Error
}

func (r *Repository[T]) GetAll(ctx context.Context) ([]T, error) {
	var entities []T
	err :=  r.database.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (r *Repository[T]) GetById(ctx context.Context, id string) (T, error) {
	var entity T
	err :=  r.database.WithContext(ctx).First(&entity, id).Error
	return entity, err
}

func (r *Repository[T]) GetByName(ctx context.Context, name string) (T, error) {
	var entity T
	err :=  r.database.WithContext(ctx).First(&entity, "name = ?", name).Error
	return entity, err
}