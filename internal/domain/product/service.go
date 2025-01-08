package product

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	domainErr "github.com/tclutin/shoppinglist-api/internal/domain/errors"
	"log/slog"
)

type Repository interface {
	Create(ctx context.Context, product Product) (uint64, error)
	Delete(ctx context.Context, productID uint64) error
	GetById(ctx context.Context, productID uint64) (Product, error)
	GetCategories(ctx context.Context) ([]Category, error)
	GetProductsByCategoryId(ctx context.Context, categoryID uint64) ([]ProductName, error)
	GetByProductNameId(ctx context.Context, productNameID uint64) (ProductName, error)
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

func (s *Service) Create(ctx context.Context, product Product) (uint64, error) {
	return s.repo.Create(ctx, product)
}

func (s *Service) RemoveProduct(ctx context.Context, productID uint64) error {
	_, err := s.repo.GetById(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainErr.ErrProductNotFound
		}

		return fmt.Errorf("failed to get product: %w", err)
	}

	return s.repo.Delete(ctx, productID)
}

func (s *Service) GetById(ctx context.Context, productID uint64) (Product, error) {
	product, err := s.repo.GetById(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return product, domainErr.ErrProductNotFound
		}

		return product, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

func (s *Service) GetByProductNameId(ctx context.Context, productNameID uint64) (ProductName, error) {
	productName, err := s.repo.GetByProductNameId(ctx, productNameID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return productName, domainErr.ErrProductNotFound
		}

		return productName, fmt.Errorf("failed to get product: %w", err)
	}

	return productName, nil
}

func (s *Service) GetCategories(ctx context.Context) ([]Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s *Service) GetProductsByCategoryId(ctx context.Context, categoryID uint64) ([]ProductName, error) {
	return s.repo.GetProductsByCategoryId(ctx, categoryID)
}
