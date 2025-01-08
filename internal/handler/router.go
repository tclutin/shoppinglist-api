package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/tclutin/shoppinglist-api/docs"
	"github.com/tclutin/shoppinglist-api/internal/config"
	"github.com/tclutin/shoppinglist-api/internal/domain"
	"github.com/tclutin/shoppinglist-api/internal/handler/auth"
	"github.com/tclutin/shoppinglist-api/internal/handler/group"
	"github.com/tclutin/shoppinglist-api/internal/handler/middleware"
	"github.com/tclutin/shoppinglist-api/internal/handler/product"
	"github.com/tclutin/shoppinglist-api/internal/handler/user"
	"github.com/tclutin/shoppinglist-api/pkg/logger"
	"net/http"
)

func NewRouter(cfg *config.Config, logger logger.Logger, services *domain.Services) *gin.Engine {
	if cfg.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	if cfg.IsDev() {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	root := router.Group("/api")
	{
		auth.NewAuthHandler(logger, services.Auth).Init(root, services.Auth)
		user.NewGroupHandler(logger, services.User).Init(root, services.Auth)
		group.NewGroupHandler(logger, services.Group).Init(root, services.Auth)
		product.NewGroupHandler(logger, services.Product).Init(root, services.Auth)
	}

	return router
}
