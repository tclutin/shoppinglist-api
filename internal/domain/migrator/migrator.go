package migrator

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"log"
)

type Migration struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Migration {
	return &Migration{
		pool: pool,
	}
}

func (m *Migration) Init(ctx context.Context) {
	m.Up()
}
func (m *Migration) Up() {
	db := stdlib.OpenDBFromPool(m.pool)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalln("error applying migrations", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalln("error applying migrations", err)
	}
}
