package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/agumiroff/BigTechProject/order/v1/internal/config/env"
)

var appConfig *config

type config struct {
	HTTPConfig      OrderConfig
	PostgressConfig PostgressConfig
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

	postgressConfig, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		HTTPConfig:      httpConfig,
		PostgressConfig: postgressConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
