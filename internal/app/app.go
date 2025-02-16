package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/shoppinglist-api/internal/config"
	"github.com/tclutin/shoppinglist-api/internal/domain"
	migrator "github.com/tclutin/shoppinglist-api/internal/domain/migrator"
	"github.com/tclutin/shoppinglist-api/internal/handler"
	"github.com/tclutin/shoppinglist-api/internal/repository"
	"github.com/tclutin/shoppinglist-api/pkg/client/postgresql"
	"github.com/tclutin/shoppinglist-api/pkg/jwt/manager"
	"github.com/tclutin/shoppinglist-api/pkg/logger"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	httpServer *http.Server
	logger     logger.Logger
	pool       *pgxpool.Pool
}

func New() *App {
	cfg := config.MustLoad()

	customLogger := logger.New(cfg.IsProd())

	dsn := fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database)

	pool := postgresql.NewPool(context.Background(), dsn)

	migrator := migrator.New(pool)
	migrator.Init(context.Background())

	tokenManager := manager.MustLoadTokenManager(cfg.JWT.Secret)

	repos := repository.NewRepositories(pool)

	services := domain.NewServices(cfg, tokenManager, repos)

	router := handler.NewRouter(cfg, customLogger, services)

	return &App{
		httpServer: &http.Server{
			Addr:           net.JoinHostPort(cfg.HTTPServer.Host, cfg.HTTPServer.Port),
			Handler:        router,
			MaxHeaderBytes: 1 << 20,
			WriteTimeout:   5 * time.Second,
			ReadTimeout:    5 * time.Second,
		},
		logger: customLogger,
		pool:   pool,
	}
}

func (a *App) Run(ctx context.Context) {
	a.logger.Info("App is starting...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				a.logger.Error("Server stopped with error", slog.Any("error", err))
				os.Exit(1)
			}
		}
	}()

	a.logger.Info("App started successfully")

	<-quit
	a.logger.Info("Received shutdown signal, stopping app...")
	a.Stop(ctx)
}

func (a *App) Stop(ctx context.Context) {
	a.logger.Info("App is shutting down...")

	a.pool.Close()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.logger.Error("Server stopped with error", slog.Any("error", err))
		os.Exit(1)
	}

	a.logger.Info("App shutdown completed")
}
