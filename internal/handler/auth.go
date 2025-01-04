package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type AuthService interface {
	SignUp(ctx context.Context)
	/*	SignIn(ctx context.Context)
		Logout(ctx context.Context)
		Me(ctx context.Context)*/
}

type AuthHandler struct {
	logger      *slog.Logger
	authService AuthService
}

func NewAuthHandler(logger *slog.Logger, service AuthService) *AuthHandler {
	return &AuthHandler{
		logger:      logger.With("authHandler"),
		authService: service,
	}
}

func (a *AuthHandler) Init(router *gin.RouterGroup) {
	authRouter := router.Group("auth")
	{
		authRouter.POST("/signup", nil)
		authRouter.POST("/login", nil)
		authRouter.GET("/who", nil)
		authRouter.POST("/refresh", nil)
	}
}
