package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
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
			pool.Close()
			slog.Error("failed to ping database, retrying...", "retry_count", i+1, "error", err)
			time.Sleep(3 * time.Second)
			continue
		}

		return pool
	}

	log.Fatalln(fmt.Errorf("failed to connect to the database after %d retries", maxRetries))
	return nil
}
