package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// mongoEnvConfig defines configuration parameters for MongoDB connection
type mongoEnvConfig struct {
	Host     string `env:"MONGO_HOST,required"`
	Port     int    `env:"MONGO_PORT,required"`
	Database string `env:"MONGO_INITDB_DATABASE,required"`
	AuthDB   string `env:"MONGO_AUTH_DB,required"`
	Username string `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password string `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
	Path     string `env:"MONGO_MIGRATIONS_DIR,required"`
}

type mongoConfig struct {
	raw mongoEnvConfig
}

// NewMongoConfig parses environment variables into a mongoConfig instance
func NewMongoConfig() (*mongoConfig, error) {
	var raw mongoEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("failed to parse Mongo env config: %w", err)
	}
	return &mongoConfig{raw: raw}, nil
}

// URI builds the MongoDB connection URI string
func (c *mongoConfig) URI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/%s?authSource=%s",
		c.raw.Username,
		c.raw.Password,
		c.raw.Host,
		c.raw.Port,
		c.raw.Database,
		c.raw.AuthDB,
	)
}

// Database returns the configured database name
func (c *mongoConfig) DBName() string {
	return c.raw.Database
}

func (c *mongoConfig) MigrationPath() string {
	return c.raw.Path
}
