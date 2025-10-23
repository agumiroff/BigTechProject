package env

import (
	"fmt"
	"net"

	"github.com/caarlos0/env/v11"
)

type grpcEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type grpcConfig struct {
	raw grpcEnvConfig
}

func NewGRPCConfig() (*grpcConfig, error) {
	var raw grpcEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("failed to parse grpc env config: %w", err)
	}
	return &grpcConfig{
		raw: raw,
	}, nil
}

func (c *grpcConfig) Address() string {
	return net.JoinHostPort(c.raw.Host, c.raw.Port)
}
