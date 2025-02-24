package persistence

import (
	"belcamp/internal/domain/entity"
	"context"
)

type CategoryRepository struct {
	*GormRepository[entity.Category]
}

// Custom repository methods
func (r *CategoryRepository) FindBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	var category entity.Category
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
