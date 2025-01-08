package domain

import (
	"github.com/tclutin/shoppinglist-api/internal/config"
	"github.com/tclutin/shoppinglist-api/internal/domain/auth"
	"github.com/tclutin/shoppinglist-api/internal/domain/group"
	"github.com/tclutin/shoppinglist-api/internal/domain/product"
	"github.com/tclutin/shoppinglist-api/internal/domain/user"
	"github.com/tclutin/shoppinglist-api/internal/repository"
	"github.com/tclutin/shoppinglist-api/pkg/jwt/manager"
)

type Services struct {
	Auth    *auth.Service
	User    *user.Service
	Group   *group.Service
	Product *product.Service
}

func NewServices(cfg *config.Config, tokenManager manager.Manager, repos *repository.Repository) *Services {
	userService := user.NewService(repos.User)
	authService := auth.NewService(cfg, userService, tokenManager, repos.Session)
	productService := product.NewService(repos.Product)
	groupService := group.NewService(repos.Group, repos.Member, productService)

	return &Services{
		Auth:    authService,
		User:    userService,
		Group:   groupService,
		Product: productService,
	}
}
