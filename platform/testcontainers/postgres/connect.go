package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func connectPostgresClient(ctx context.Context, uri string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, errors.Errorf("failed to create new pool %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, errors.Errorf("failed to ping postgres %v", err)
	}

	return pool, nil
}
