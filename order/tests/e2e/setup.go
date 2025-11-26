package integration

import (
	"context"
	"os"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	"github.com/agumiroff/BigTechProject/platform/testcontainers"
	"github.com/agumiroff/BigTechProject/platform/testcontainers/app"
	network "github.com/agumiroff/BigTechProject/platform/testcontainers/network"
	"github.com/agumiroff/BigTechProject/platform/testcontainers/path"
	"github.com/agumiroff/BigTechProject/platform/testcontainers/postgres"
)

// TestEnvironment — структура для хранения ресурсов тестового окружения
type TestEnvironment struct {
	Network  *network.Network
	Postgres *postgres.PostgresContainer
	App      *app.Container
}

// setupTestEnvironment — подготавливает тестовое окружение: сеть, контейнеры и возвращает структуру с ресурсами
func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	// Шаг 1: Создаём общую Docker-сеть
	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "✅ Сеть успешно создана")

	// Получаем переменные окружения для PostgreSQL с проверкой на наличие
	postgresUsername := getEnvWithLogging(ctx, testcontainers.PostgresUserKey)
	postgresPassword := getEnvWithLogging(ctx, testcontainers.PostgresPasswordKey)
	postgresImageName := getEnvWithLogging(ctx, testcontainers.PostgresImageNameKey)
	postgresDatabase := getEnvWithLogging(ctx, testcontainers.PostgresDatabaseKey)

	// Получаем порт gRPC для waitStrategy
	grpcPort := getEnvWithLogging(ctx, grpcPortKey)
	grpcHost := getEnvWithLogging(ctx, grpcHostKey)

	// Шаг 2: Запускаем контейнер с PostgreSQL
	generatedPostgres, err := postgres.NewPostgresContainer(ctx,
		postgres.WithNetworkName(generatedNetwork.Name()),
		postgres.WithContainerName(testcontainers.PostgresContainerName),
		postgres.WithImageName(postgresImageName),
		postgres.WithDatabase(postgresDatabase),
		postgres.WithAuth(postgresUsername, postgresPassword),
		postgres.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер PostgreSQL", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер PostgreSQL успешно запущен")

	// Шаг 3: Запускаем контейнер с приложением
	projectRoot := path.GetProjectRoot()

	// Переопределяем хост PostgreSQL для подключения к контейнеру из testcontainers
	appEnv := map[string]string{
		testcontainers.PostgresHostKey:       generatedPostgres.Config().ContainerName,
		testcontainers.PostgresPortKey:       testcontainers.PostgresDefaultPort,
		testcontainers.PostgresDatabaseKey:   postgresDatabase,
		testcontainers.PostgresUserKey:       postgresUsername,
		testcontainers.PostgresPasswordKey:   postgresPassword,
		testcontainers.PostgresMigrationPath: "/app/migrations",
		grpcPortKey:                          grpcPort,
		testcontainers.MongoLoggerLevel:      loggerLevelValue,
		grpcHostKey:                          grpcHost,
	}

	// Создаем настраиваемую стратегию ожидания с увеличенным таймаутом
	waitStrategy := wait.ForListeningPort(nat.Port(grpcPort + "/tcp")).
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(AppName),
		app.WithPort(grpcPort),
		app.WithDockerfile(projectRoot, Dockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения успешно запущен")

	logger.Info(ctx, "🎉 Тестовое окружение готово")
	return &TestEnvironment{
		Network:  generatedNetwork,
		Postgres: generatedPostgres,
		App:      appContainer,
	}
}

// getEnvWithLogging возвращает значение переменной окружения с логированием
func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "Переменная окружения не установлена", zap.String("key", key))
	}

	return value
}
