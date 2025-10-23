package config

import (
	"log"
	"os"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	GRPCConfig  InventoryConfig
	MongoConfig MongoConfig
}

func Load() error {
	err := godotenv.Load()
	if err != nil && os.IsNotExist(err) {
		log.Printf("Error loading config %s", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		return err
	}

	mongoConfig, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		GRPCConfig:  grpcConfig,
		MongoConfig: mongoConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
