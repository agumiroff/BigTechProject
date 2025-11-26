package integration

import "time"

const (
	// projectName - имя проекта для Docker-контейнеров и сети
	projectName = "inventory-service"

	// partsCollectionName - имя коллекции MongoDB для запчастей
	partsCollectionName = "inventory"

	// Параметры для контейнеров
	AppName    = "inventory-app"
	Dockerfile = "inventory/Dockerfile"

	// Переменные окружения приложения
	grpcPortKey = "GRPC_PORT"
	grpcHostKey = "GRPC_HOST"

	// Значения переменных окружения
	loggerLevelValue = "debug"
	startupTimeout   = 3 * time.Minute
)
