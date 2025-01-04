package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tclutin/shoppinglist-api/internal/config"
	"log/slog"
	"net/http"
)

type Handler struct {
	logger *slog.Logger
	cfg    *config.Config
}

func New(logger *slog.Logger, cfg *config.Config) *Handler {
	return &Handler{
		logger: logger.With("handler", "handler"),
		cfg:    cfg,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	if h.cfg.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	if h.cfg.IsDev() {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	root := router.Group("/api")
	{
		NewAuthHandler(h.logger, nil).Init(root)
	}

	return router
}
