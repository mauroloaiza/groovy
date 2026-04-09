package db

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	url := os.Getenv("DB_URL")
	if url == "" {
		url = "postgres://groovy:groovy@localhost:5432/groovy?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("db connect: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	if err := migrate(ctx, pool); err != nil {
		return nil, fmt.Errorf("db migrate: %w", err)
	}

	return pool, nil
}

func migrate(ctx context.Context, pool *pgxpool.Pool) error {
	entries, err := migrations.ReadDir("migrations")
	if err != nil {
		return err
	}

	for _, e := range entries {
		sql, err := migrations.ReadFile("migrations/" + e.Name())
		if err != nil {
			return err
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("migration %s: %w", e.Name(), err)
		}
	}

	return nil
}
