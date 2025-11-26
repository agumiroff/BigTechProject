package integration

import "time"

const (
	// projectName - имя проекта для Docker-контейнеров и сети
	projectName = "order-service"

	// ordersTableName - имя таблицы PostgreSQL для заказов
	ordersTableName = "orders"

	// orderPartsTableName - имя таблицы PostgreSQL для частей заказа
	orderPartsTableName = "order"

	// Параметры для контейнеров
	AppName    = "order-app"
	Dockerfile = "order/Dockerfile"

	// Переменные окружения приложения
	grpcPortKey = "GRPC_PORT"
	grpcHostKey = "GRPC_HOST"

	// Значения переменных окружения
	loggerLevelValue = "debug"
	startupTimeout   = 3 * time.Minute

	dataBaseName = "POSTGRES_DATABASE"
)
