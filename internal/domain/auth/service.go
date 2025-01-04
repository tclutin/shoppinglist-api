package auth

import (
	"context"
	"github.com/tclutin/shoppinglist-api/internal/domain/user"
	"github.com/tclutin/shoppinglist-api/pkg/jwt/manager"
	"log/slog"
)

type UserService interface {
	GetById(ctx context.Context, userID uint64) (user.User, error)
	Create(ctx context.Context, user user.User) (uint64, error)
}

type Service struct {
	logger       *slog.Logger
	tokenManager manager.TokenManager
	userService  UserService
}

func NewService(logger *slog.Logger, userService UserService, tokenManager manager.TokenManager) *Service {
	return &Service{
		logger:       logger.With("service", "user_service"),
		tokenManager: tokenManager,
		userService:  userService,
	}
}

func (s *Service) LogIn(ctx context.Context, dto LogInDTO) (user.User, error) {
	panic("implement me")
}

func (s *Service) SignUp(ctx context.Context, dto SignUpDTO) (user.User, error) {
	panic("implement me")
}

func (s *Service) Who(ctx context.Context, userID uint64) (user.User, error) {
	usr, err := s.userService.GetById(ctx, userID)
	if err != nil {
		return usr, err
	}

	return usr, nil
}
