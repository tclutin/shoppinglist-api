package auth

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tclutin/shoppinglist-api/internal/domain/auth"
	domainErr "github.com/tclutin/shoppinglist-api/internal/domain/errors"
	"github.com/tclutin/shoppinglist-api/internal/domain/user"
	mw "github.com/tclutin/shoppinglist-api/internal/handler/middleware"
	"github.com/tclutin/shoppinglist-api/pkg/logger"
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
	logger  logger.Logger
	service Service
}

func NewAuthHandler(logger logger.Logger, service Service) *Handler {
	return &Handler{
		logger:  logger.With("handler", "auth_handler"),
		service: service,
	}
}

func (h *Handler) Init(router *gin.RouterGroup, authService *auth.Service) {
	authRouter := router.Group("auth")
	{
		authRouter.POST("/signup", h.SignUp)
		authRouter.POST("/login", h.LogIn)
		authRouter.POST("/refresh", h.Refresh)
		authRouter.GET("/who", mw.AuthMiddleware(authService), h.Who)
	}
}

// @Summary		SignUp
// @Description	Create new user
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body		SignUpRequest	true	"Create new user"
// @Success		201		{object}	TokenResponse
// @Failure		422		{object}	response.APIError
// @Failure		409		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/auth/signup [post]
func (h *Handler) SignUp(c *gin.Context) {
	var request SignUpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			response.NewAPIError(http.StatusUnprocessableEntity, err.Error(), nil))
		return
	}

	tokens, err := h.service.SignUp(c.Request.Context(), auth.SignUpDTO{
		Username: request.Username,
		Password: request.Password,
		Gender:   request.Gender,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrUserAlreadyExists) {
			c.AbortWithStatusJSON(
				http.StatusConflict,
				response.NewAPIError(http.StatusConflict, err.Error(), nil))
			return
		}

		h.logger.Error("error occurred while processing SignUp", slog.Any("error", err))
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError(http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusCreated, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// @Summary		LogIn
// @Description	Log in your account
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body		LogInRequest	true	"Log in your account"
// @Success		200		{object}	TokenResponse
// @Failure		422		{object}	response.APIError
// @Failure		400		{object}	response.APIError
// @Failure		404		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/auth/login [post]
func (h *Handler) LogIn(c *gin.Context) {
	var request LogInRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			response.NewAPIError(http.StatusUnprocessableEntity, err.Error(), nil))
		return
	}

	tokens, err := h.service.LogIn(c.Request.Context(), auth.LogInDTO{
		Username: request.Username,
		Password: request.Password,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound,
				response.NewAPIError(http.StatusNotFound, err.Error(), nil))
			return
		}

		if errors.Is(err, domainErr.ErrUserNotValid) {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				response.NewAPIError(http.StatusBadRequest, err.Error(), nil))
			return
		}

		h.logger.Error("error occurred while processing LogIn", slog.Any("error", err))
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError(http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// @Summary		Refresh
// @Description	Refresh your token
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body		RefreshTokenRequest	true	"Refresh your token"
// @Success		200		{object}	TokenResponse
// @Failure		422		{object}	response.APIError
// @Failure		400		{object}	response.APIError
// @Failure		404		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var request RefreshTokenRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			response.NewAPIError(http.StatusUnprocessableEntity, err.Error(), nil))
		return
	}

	uuid, err := uuid.Parse(request.RefreshToken)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			response.NewAPIError(http.StatusBadRequest, "failed to parse refresh token", nil))
		return
	}

	tokens, err := h.service.Refresh(c.Request.Context(), auth.RefreshTokenDTO{
		RefreshToken: uuid,
	})

	if err != nil {
		if errors.Is(err, domainErr.ErrSessionNotFound) {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				response.NewAPIError(http.StatusNotFound, err.Error(), nil))
			return
		}

		if errors.Is(err, domainErr.ErrRefreshTokenExpired) {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				response.NewAPIError(http.StatusBadRequest, err.Error(), nil))
			return
		}
		h.logger.Error("error occurred while processing Refresh", slog.Any("error", err))
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError(http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// @Security		ApiKeyAuth
// @Summary		Who
// @Description	Get current user
// @Tags			auth
// @Accept			json
// @Produce		json
// @Success		200		{object}	CurrentUserResponse
// @Failure		401		{object}	response.APIError
// @Failure		404		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/auth/who [get]
func (h *Handler) Who(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			response.NewAPIError(http.StatusUnauthorized, domainErr.ErrMissingCredentials.Error(), nil))
		return
	}

	usr, err := h.service.Who(c.Request.Context(), userID.(uint64))
	if err != nil {
		if errors.Is(err, domainErr.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound,
				response.NewAPIError(http.StatusNotFound, err.Error(), nil))
			return
		}

		h.logger.Error("error occurred while processing Who", slog.Any("error", err))
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError(http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusOK, CurrentUserResponse{
		UserID:    usr.UserID,
		Username:  usr.Username,
		Gender:    usr.Gender,
		CreatedAt: usr.CreatedAt,
	})
}
