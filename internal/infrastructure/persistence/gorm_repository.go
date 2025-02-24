package persistence

import (
	"belcamp/internal/domain/repository"
	"belcamp/internal/infrastructure/errors"
	"context"

	"gorm.io/gorm"
)

type GormRepository[T any] struct {
	db *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) repository.Repository[T] {
	return &GormRepository[T]{db: db}
}

func (r *GormRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *GormRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return &entity, nil
}

func (r *GormRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *GormRepository[T]) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(new(T), id).Error
}

func (r *GormRepository[T]) List(ctx context.Context, page, pageSize int) ([]T, int64, error) {
	var entities []T
	var total int64

	if err := r.db.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
