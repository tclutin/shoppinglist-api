package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

const (
	maxRetries = 3
)

type Client interface {
}

func NewPool(ctx context.Context, dsn string) *pgxpool.Pool {
	for i := 0; i < maxRetries; i++ {
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			slog.Error("failed to connect to the database", "retry_count", i+1, "error", err)
			time.Sleep(3 * time.Second)
			continue
		}

		if err = pool.Ping(ctx); err != nil {
			slog.Error("failed to ping to the database", "retry_count", i+1, "error", err)
			time.Sleep(3 * time.Second)
			continue
		}

		return pool
	}

	return nil
}
