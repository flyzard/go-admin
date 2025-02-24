package persistence

import (
	"belcamp/internal/domain/entity"
	"context"
)

type ProductRepository struct {
	*GormRepository[entity.Product]
}

// Custom repository methods
func (r *ProductRepository) FindBySlug(ctx context.Context, slug string) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindByCategory(ctx context.Context, categoryID uint) ([]entity.Product, error) {
	var products []entity.Product
	if err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
