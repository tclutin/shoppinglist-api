package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/shoppinglist-api/internal/domain/auth"
	domainErr "github.com/tclutin/shoppinglist-api/internal/domain/errors"
	"github.com/tclutin/shoppinglist-api/internal/domain/group"
	mw "github.com/tclutin/shoppinglist-api/internal/handler/middleware"
	"github.com/tclutin/shoppinglist-api/pkg/logger"
	"github.com/tclutin/shoppinglist-api/pkg/response"
	"log/slog"
	"net/http"
)

type Service interface {
	GetGroupsByUserId(ctx context.Context, userId uint64) ([]group.GroupDTO, error)
}

type Handler struct {
	logger  logger.Logger
	service Service
}

func NewGroupHandler(logger logger.Logger, service Service) *Handler {
	return &Handler{
		logger:  logger.With("handler", "user_handler"),
		service: service,
	}
}

func (h *Handler) Init(router *gin.RouterGroup, authService *auth.Service) {
	usersRouter := router.Group("users", mw.AuthMiddleware(authService))
	{
		usersRouter.GET("/groups", h.GetUserGroups)
	}
}

// @Security		ApiKeyAuth
// @Summary		GetUserGroups
// @Description	get user groups
// @Tags			users
// @Accept			json
// @Produce		json
// @Success		200		{object}	group.GroupUserDTO
// @Failure		401		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/users/groups [get]
func (h *Handler) GetUserGroups(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			response.NewAPIError(http.StatusUnauthorized, domainErr.ErrMissingCredentials.Error(), nil))
		return
	}

	groups, err := h.service.GetGroupsByUserId(c.Request.Context(), userID.(uint64))
	if err != nil {
		h.logger.Error("error occurred while processing GetGroupsByUserId", slog.Any("error", err))
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError(http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusOK, groups)
}
