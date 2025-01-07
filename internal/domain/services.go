package domain

import (
	"github.com/tclutin/shoppinglist-api/internal/config"
	"github.com/tclutin/shoppinglist-api/internal/domain/auth"
	"github.com/tclutin/shoppinglist-api/internal/domain/group"
	"github.com/tclutin/shoppinglist-api/internal/domain/product"
	"github.com/tclutin/shoppinglist-api/internal/domain/user"
	"github.com/tclutin/shoppinglist-api/internal/repository"
	"github.com/tclutin/shoppinglist-api/pkg/jwt/manager"
	"log/slog"
)

type Services struct {
	Auth    *auth.Service
	User    *user.Service
	Group   *group.Service
	Product *product.Service
}

func NewServices(logger *slog.Logger, cfg *config.Config, tokenManager manager.Manager, repos *repository.Repository) *Services {
	userService := user.NewService(logger, repos.User)
	authService := auth.NewService(logger, cfg, userService, tokenManager, repos.Session)
	productService := product.NewService(repos.Product, logger)
	groupService := group.NewService(repos.Group, repos.Member, productService, logger)

	return &Services{
		Auth:    authService,
		User:    userService,
		Group:   groupService,
		Product: productService,
	}
}
