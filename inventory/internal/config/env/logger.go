package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type loggerEnvConfig struct {
	Level  string `env:"LOG_LEVEL" envDefault:"info"`
	AsJSON bool   `env:"LOG_AS_JSON" envDefault:"false"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("failed to parse logger env config: %w", err)
	}
	return &loggerConfig{
		raw: raw,
	}, nil
}

func (c *loggerConfig) Level() string {
	return c.raw.Level
}

func (c *loggerConfig) AsJson() bool {
	return c.raw.AsJSON
}
