package service

import (
	"belcamp/internal/domain/repository"
	"belcamp/internal/domain/valueobject"
	"context"
)

type CRUDService[T any] struct {
	repo repository.Repository[T]
}

func NewCRUDService[T any](repo repository.Repository[T]) *CRUDService[T] {
	return &CRUDService[T]{repo: repo}
}

func (s *CRUDService[T]) Create(ctx context.Context, entity *T) error {
	return s.repo.Create(ctx, entity)
}

func (s *CRUDService[T]) Get(ctx context.Context, id uint) (*T, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CRUDService[T]) Update(ctx context.Context, entity *T) error {
	return s.repo.Update(ctx, entity)
}

func (s *CRUDService[T]) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *CRUDService[T]) List(ctx context.Context, pagination *valueobject.Pagination) ([]T, *valueobject.Pagination, error) {
	entities, total, err := s.repo.List(ctx, pagination.Page, pagination.PageSize)
	if err != nil {
		return nil, nil, err
	}
	pagination.Total = total
	return entities, pagination, nil
}
