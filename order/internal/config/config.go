package config

import (
	"log"
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
	if err != nil && os.IsNotExist(err) {
		log.Printf("Error loading config %s", err)
	}

	httpConfig, err := env.NewHTTPConfig()
	if err != nil {
		log.Printf("Error loading HTTP config: %v", err)
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
