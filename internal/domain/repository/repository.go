// Package repository provides the repository interface for the application.
package repository

import "context"

type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id uint) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page, pageSize int, preloadFields ...string) ([]T, int64, error)
}
