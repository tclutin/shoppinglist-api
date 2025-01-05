package domain

import (
	"github.com/tclutin/shoppinglist-api/internal/config"
	"github.com/tclutin/shoppinglist-api/internal/domain/auth"
	"github.com/tclutin/shoppinglist-api/internal/domain/user"
	"github.com/tclutin/shoppinglist-api/internal/repository"
	"github.com/tclutin/shoppinglist-api/pkg/jwt/manager"
	"log/slog"
)

type Services struct {
	Auth *auth.Service
	User *user.Service
}

func NewServices(logger *slog.Logger, cfg *config.Config, tokenManager manager.Manager, repos *repository.Repository) *Services {
	userService := user.NewService(logger, repos.User)
	authService := auth.NewService(logger, cfg, userService, tokenManager, repos.Session)

	return &Services{
		Auth: authService,
		User: userService,
	}
}
