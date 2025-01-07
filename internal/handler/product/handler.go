package product

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/shoppinglist-api/internal/domain/auth"
	"github.com/tclutin/shoppinglist-api/internal/domain/product"
	mw "github.com/tclutin/shoppinglist-api/internal/handler/middleware"
	"github.com/tclutin/shoppinglist-api/pkg/response"
	"log/slog"
	"net/http"
	"strconv"
)

type Service interface {
	GetProductsByCategoryId(ctx context.Context, categoryID uint64) ([]product.ProductName, error)
	GetCategories(ctx context.Context) ([]product.Category, error)
}

type Handler struct {
	logger  *slog.Logger
	service Service
}

func NewGroupHandler(logger *slog.Logger, service Service) *Handler {
	return &Handler{
		logger:  logger.With("handler", "productHandler"),
		service: service,
	}
}

func (h *Handler) Init(router *gin.RouterGroup, authService *auth.Service) {
	productGroup := router.Group("/products", mw.AuthMiddleware(authService))
	{
		productGroup.GET("/categories", h.GetCategories)
		productGroup.GET("/:category_id", h.GetProductsByCategoryId)
	}
}

// @Security		ApiKeyAuth
// @Summary		GetCategories
// @Description	get categories
// @Tags			products
// @Accept			json
// @Produce		json
// @Success		200		{object}	product.Category
// @Failure		401		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/products/categories [get]
func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError(http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusOK, categories)
}

// @Security		ApiKeyAuth
// @Summary		GetProductsByCategory
// @Description	get products by category id
// @Tags			products
// @Accept			json
// @Produce		json
// @Param			category_id	path		string	true	"category_id"
// @Success		200		{object}	product.ProductName
// @Failure		401		{object}	response.APIError
// @Failure		422		{object}	response.APIError
// @Failure		500		{object}	response.APIError
// @Router			/products/{category_id} [get]
func (h *Handler) GetProductsByCategoryId(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("category_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			response.NewAPIError(http.StatusUnprocessableEntity, "':category_id' is not correct", nil))
		return
	}

	products, err := h.service.GetProductsByCategoryId(c.Request.Context(), categoryID)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.NewAPIError(http.StatusInternalServerError, "Internal server error", nil))
		return
	}

	c.JSON(http.StatusOK, products)
}
