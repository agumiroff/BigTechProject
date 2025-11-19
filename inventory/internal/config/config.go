package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/config/env"
)

var appConfig *config

type config struct {
	GRPC   InventoryConfig
	Mongo  MongoConfig
	Logger LoggerConfig
}

func Load() error {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		// Log actual errors, ignore missing .env
		return err
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		return err
	}

	mongoConfig, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		GRPC:   grpcConfig,
		Mongo:  mongoConfig,
		Logger: loggerConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
