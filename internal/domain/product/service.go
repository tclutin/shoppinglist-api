package product

import (
	"context"
	"log/slog"
)

type Repository interface {
	Create(ctx context.Context, product Product) (uint64, error)
	GetCategories(ctx context.Context) ([]Category, error)
	GetProductsByCategoryId(ctx context.Context, categoryID uint64) ([]ProductName, error)
}

type Service struct {
	logger *slog.Logger
	repo   Repository
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		logger: logger.With("service", "product_service"),
		repo:   repo,
	}
}

func (s *Service) GetCategories(ctx context.Context) ([]Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s *Service) GetProductsByCategoryId(ctx context.Context, categoryID uint64) ([]ProductName, error) {
	return s.repo.GetProductsByCategoryId(ctx, categoryID)
}
