package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host     string `env:"POSTGRES_HOST,required"`
	Port     int    `env:"POSTGRES_PORT,required"`
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	DBName   string `env:"POSTGRES_DB,required"`
	SSLMode  string `env:"POSTGRES_SSLMODE" envDefault:"disable"`
	MigPath  string `env:"POSTGRES_MIGRATION_PATH", required`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

// NewPostgresConfig parses environment variables into a postgresConfig instance
func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("failed to parse Postgres env config: %w", err)
	}
	return &postgresConfig{raw: raw}, nil
}

// DSN builds the Postgres connection string
func (c *postgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.raw.Host,
		c.raw.Port,
		c.raw.User,
		c.raw.Password,
		c.raw.DBName,
		c.raw.SSLMode,
	)
}

// DBName returns the configured database name
func (c *postgresConfig) DBName() string {
	return c.raw.DBName
}

func (c *postgresConfig) MigPath() string {
	return c.raw.MigPath
}
