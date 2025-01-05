package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	domainErr "github.com/tclutin/shoppinglist-api/internal/domain/errors"
	"log/slog"
)

type Repository interface {
	Create(ctx context.Context, user User) (uint64, error)
	GetById(ctx context.Context, userID uint64) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
}

type Service struct {
	logger *slog.Logger
	repo   Repository
}

func NewService(logger *slog.Logger, repo Repository) *Service {
	return &Service{
		logger: logger.With("service", "UserService"),
		repo:   repo,
	}
}

func (s *Service) Create(ctx context.Context, user User) (uint64, error) {
	userID, err := s.repo.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

func (s *Service) GetById(ctx context.Context, userID uint64) (User, error) {
	user, err := s.repo.GetById(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, domainErr.ErrUserNotFound
		}

		return user, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *Service) GetByUsername(ctx context.Context, username string) (User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, domainErr.ErrUserNotFound
		}

		return user, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
