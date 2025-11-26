package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
)

const (
	postgresStartupTimeout = 1 * time.Minute
)

type PostgresContainer struct {
	container testcontainers.Container
	pool      *pgxpool.Pool
	config    *Config
}

func (c *PostgresContainer) Pool() *pgxpool.Pool {
	return c.pool
}

func (c *PostgresContainer) Config() *Config {
	return c.config
}

func (c *PostgresContainer) Terminate(ctx context.Context) error {
	if c.pool != nil {
		c.pool.Close()
	}
	if c.container != nil {
		return c.container.Terminate(ctx)
	}
	return nil
}

func NewPostgresContainer(ctx context.Context, opts ...Option) (*PostgresContainer, error) {
	cfg := buildConfig(opts...)

	container, err := startPostgresContainer(ctx, cfg)
	if err != nil {
		return nil, err
	}

	success := false
	defer func() {
		if !success {
			if err := container.Terminate(ctx); err != nil {
				cfg.Logger.Error(ctx, "failed to terminate mongo container", zap.Error(err))
			}
		}
	}()

	cfg.Host, cfg.Port, err = getContainerHostPort(ctx, container)
	if err != nil {
		return nil, err
	}

	uri := buildPostgresURI(cfg)
	client, err := connectPostgresClient(ctx, uri)

	cfg.Logger.Info(ctx, "Postgres container started", zap.String("uri", uri))
	success = true

	return &PostgresContainer{
		container: container,
		pool:      client,
		config:    cfg,
	}, nil
}
