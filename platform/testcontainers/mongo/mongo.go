package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	mongoPort           = "27017"
	mongoStartupTimeout = 1 * time.Minute
)

type MongoContainer struct {
	container testcontainers.Container
	client    *mongo.Client
	cfg       *Config
}

func NewMongoContainer(ctx context.Context, opts ...Option) (*MongoContainer, error) {
	cfg := buildConfig(opts...)

	container, err := startMongoContainer(ctx, cfg)
	if err != nil {
		return nil, errors.Errorf("error starting mongo container %v", err)
	}

	success := false
	defer func() {
		if !success {
			if err = container.Terminate(ctx); err != nil {
				cfg.Logger.Error(ctx, "failed to terminate mongo container", zap.Error(err))
			}
		}
	}()

	cfg.Host, cfg.Port, err = getContainerHostPort(ctx, container)
	if err != nil {
		return nil, err
	}

	uri := buildMongoURI(cfg)

	client, err := connectMongoClient(ctx, uri)
	if err != nil {
		return nil, err
	}

	cfg.Logger.Info(ctx, "Mongo container started", zap.String("uri", uri))
	success = true

	return &MongoContainer{
		container: container,
		client:    client,
		cfg:       cfg,
	}, nil
}

func (c *MongoContainer) Client() *mongo.Client {
	return c.client
}

func (c *MongoContainer) Config() *Config {
	return c.cfg
}

func (c *MongoContainer) Terminate(ctx context.Context) error {
	if err := c.client.Disconnect(ctx); err != nil {
		c.cfg.Logger.Error(ctx, "failed to disconnect mongo client", zap.Error(err))
	}

	if err := c.container.Terminate(ctx); err != nil {
		c.cfg.Logger.Error(ctx, "failed to terminate mongo container", zap.Error(err))
	}

	c.cfg.Logger.Info(ctx, "Mongo container terminated")

	return nil
}
