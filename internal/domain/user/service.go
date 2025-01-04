package user

import (
	"context"
	"log/slog"
	"os/user"
)

type Repository interface {
	Create(ctx context.Context, user User) (uint64, error)
	GetById(ctx context.Context, userID uint64) (User, error)
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

func (s *Service) Create(ctx context.Context, user user.User) (uint64, error) {
	panic("implement me")
}

func (s *Service) GetById(ctx context.Context, userID uint64) (User, error) {
	panic("implement me")
}
