package app

import (
	"io"

	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Config struct {
	Name          string
	Port          string
	Host          string
	Logger        Logger
	DockerfileDir string
	Dockerfile    string
	Env           map[string]string
	Networks      []string
	LogOutput     io.Writer
	StartupWait   wait.Strategy
}

func buildConfig(opts ...Option) *Config {
	cfg := &Config{
		Name:          defaultAppName,
		Host:          "0.0.0.0",
		Port:          defaultPort,
		Logger:        &logger.NoopLogger{},
		DockerfileDir: ".",
		Dockerfile:    "Dockerfile",
		LogOutput:     io.Discard,
		StartupWait:   wait.ForListeningPort(defaultPort + "/tcp").WithStartupTimeout(defaultStartupTimeout),
		Env:           make(map[string]string),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
