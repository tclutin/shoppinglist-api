package user

import (
	"github.com/gin-gonic/gin"
	"github.com/tclutin/shoppinglist-api/internal/domain/auth"
	mw "github.com/tclutin/shoppinglist-api/internal/handler/middleware"
	"log/slog"
)

type Service interface {
}

type Handler struct {
	logger  *slog.Logger
	service Service
}

func NewGroupHandler(logger *slog.Logger, service Service) *Handler {
	return &Handler{
		logger:  logger.With("handler", "userHandler"),
		service: service,
	}
}

func (h *Handler) Init(router *gin.RouterGroup, authService *auth.Service) {
	usersRouter := router.Group("users", mw.AuthMiddleware(authService))
	{
		usersRouter.GET("/groups", nil)
	}
}

func (h *Handler) GetUserGroups(c *gin.Context) {}
