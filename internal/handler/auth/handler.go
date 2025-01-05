package auth

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/shoppinglist-api/internal/domain/auth"
	domainErr "github.com/tclutin/shoppinglist-api/internal/domain/errors"
	"github.com/tclutin/shoppinglist-api/internal/domain/user"
	mw "github.com/tclutin/shoppinglist-api/internal/handler/middleware"
	"github.com/tclutin/shoppinglist-api/pkg/response"
	"log/slog"
	"net/http"
)

type Service interface {
	LogIn(ctx context.Context, dto auth.LogInDTO) (auth.TokenDTO, error)
	SignUp(ctx context.Context, dto auth.SignUpDTO) (auth.TokenDTO, error)
	Refresh(ctx context.Context, dto auth.RefreshTokenDTO) (auth.TokenDTO, error)
	Who(ctx context.Context, userID uint64) (user.User, error)
}

type Handler struct {
	logger  *slog.Logger
	service Service
}

func NewAuthHandler(logger *slog.Logger, service Service) *Handler {
	return &Handler{
		logger:  logger.With("authHandler"),
		service: service,
	}
}

func (a *Handler) Init(router *gin.RouterGroup, authService *auth.Service) {
	authRouter := router.Group("auth")
	{
		authRouter.POST("/signup", a.SignUp)
		authRouter.POST("/login", a.LogIn)
		authRouter.POST("/refresh", a.Refresh)
		authRouter.GET("/who", mw.AuthMiddleware(authService), a.Who)
	}
}

func (a *Handler) LogIn(c *gin.Context) {
	var request LogInRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			response.NewAPIError[string](http.StatusUnprocessableEntity, err.Error(), nil))
		return
	}

	tokens, err := a.service.LogIn(c.Request.Context(), auth.LogInDTO{
		Username: request.Username,
		Password: request.Password,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound,
				response.NewAPIError[string](http.StatusNotFound, err.Error(), nil))
			return
		}

		if errors.Is(err, domainErr.ErrUserNotValid) {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				response.NewAPIError[string](http.StatusBadRequest, err.Error(), nil))
			return
		}

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError[string](http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (a *Handler) SignUp(c *gin.Context) {
	var request SignUpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			response.NewAPIError[string](http.StatusUnprocessableEntity, err.Error(), nil))
		return
	}

	tokens, err := a.service.SignUp(c.Request.Context(), auth.SignUpDTO{
		Username: request.Username,
		Password: request.Password,
		Gender:   request.Gender,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrUserAlreadyExists) {
			c.AbortWithStatusJSON(
				http.StatusConflict,
				response.NewAPIError[string](http.StatusConflict, err.Error(), nil))
			return
		}

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError[string](http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusCreated, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (a *Handler) Who(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			response.NewAPIError[string](http.StatusUnauthorized, domainErr.ErrMissingCredentials.Error(), nil))
		return
	}

	usr, err := a.service.Who(c.Request.Context(), userID.(uint64))
	if err != nil {
		if errors.Is(err, domainErr.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound,
				response.NewAPIError[string](http.StatusNotFound, err.Error(), nil))
			return
		}

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError[string](http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusOK, CurrentUserResponse{
		UserID:    usr.UserID,
		Username:  usr.Username,
		Gender:    usr.Gender,
		CreatedAt: usr.CreatedAt,
	})
}

func (a *Handler) Refresh(c *gin.Context) {

}
