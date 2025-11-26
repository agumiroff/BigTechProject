package postgres

import (
	"context"
	"fmt"

	constants "github.com/agumiroff/BigTechProject/platform/testcontainers"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startPostgresContainer(ctx context.Context, cfg *Config) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Name:     cfg.ContainerName,
		Image:    cfg.ImageName,
		Networks: []string{cfg.NetworkName},
		Env: map[string]string{
			constants.PostgresContainerName: cfg.Username,
			constants.PostgresPasswordKey:   cfg.Password,
		},
		WaitingFor:         wait.ForListeningPort(constants.PostgresPortKey + "/tcp").WithStartupTimeout(postgresStartupTimeout),
		HostConfigModifier: defaultHostConfig(),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, errors.Errorf("failed to start mongo container: %v", err)
	}

	return container, nil
}

func getContainerHostPort(ctx context.Context, container testcontainers.Container) (string, string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", "", errors.Errorf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, constants.PostgresPortKey+"/tcp")
	if err != nil {
		return "", "", errors.Errorf("failed to get mapped port: %v", err)
	}

	return host, port.Port(), nil
}

func buildPostgresURI(cfg *Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
}
