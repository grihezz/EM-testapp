package migrations

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"tt/internal/migrations/psql"
)

func UpMigration(ctx context.Context, pool *pgxpool.Pool) error {
	goose.SetBaseFS(psql.Embed)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	db := stdlib.OpenDBFromPool(pool)

	if err := goose.Up(db, "."); err != nil {
		panic(err)
	}

	return nil
}
