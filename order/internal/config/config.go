package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/agumiroff/BigTechProject/order/v1/internal/config/env"
)

var appConfig *config

type config struct {
	HTTP     OrderConfig
	Postgres PostgressConfig
	Logger   LoggerConfig
}

func Load() error {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		// Log actual errors, ignore missing .env
		return err
	}

	httpConfig, err := env.NewHTTPConfig()
	if err != nil {
		return err
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		HTTP:     httpConfig,
		Postgres: postgresConfig,
		Logger:   loggerConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
