package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tclutin/shoppinglist-api/internal/config"
	"github.com/tclutin/shoppinglist-api/internal/domain"
	"github.com/tclutin/shoppinglist-api/internal/handler/auth"
	"log/slog"
	"net/http"
)

func NewRouter(cfg *config.Config, logger *slog.Logger, services *domain.Services) *gin.Engine {
	if cfg.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	if cfg.IsDev() {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	root := router.Group("/api")
	{
		auth.NewAuthHandler(logger, services.Auth).Init(root, services.Auth)
	}

	return router
}
