package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

const (
	maxRetries = 3
)

type Client interface {
}

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	for i := 0; i < maxRetries; i++ {
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			slog.Error("failed to connect to the database", "retry_count", i+1, "error", err)
			time.Sleep(3 * time.Second)
			continue
		}

		return pool, nil
	}

	return nil, fmt.Errorf("failed to connect to the database after %d retries", maxRetries)
}
