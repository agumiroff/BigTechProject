package integration

import (
	"context"
	"log"
	"os"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	"github.com/agumiroff/BigTechProject/platform/testcontainers"
	"github.com/agumiroff/BigTechProject/platform/testcontainers/app"
	"github.com/agumiroff/BigTechProject/platform/testcontainers/mongo"
	"github.com/agumiroff/BigTechProject/platform/testcontainers/network"
	"github.com/agumiroff/BigTechProject/platform/testcontainers/path"
)

// TestEnvironment — структура для хранения ресурсов тестового окружения
type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.MongoContainer
	App     *app.Container
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

	// Получаем переменные окружения для MongoDB с проверкой на наличие
	mongoUsername := getEnvWithLogging(ctx, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogging(ctx, testcontainers.MongoPasswordKey)
	mongoImageName := getEnvWithLogging(ctx, testcontainers.MongoImageNameKey)
	mongoDatabase := getEnvWithLogging(ctx, testcontainers.MongoDatabaseKey)

	// Получаем порт gRPC для waitStrategy
	grpcPort := getEnvWithLogging(ctx, grpcPortKey)
	grpcHost := getEnvWithLogging(ctx, grpcHostKey)
	log.Printf("=================== %s", grpcHost)

	// Шаг 2: Запускаем контейнер с MongoDB
	generatedMongo, err := mongo.NewMongoContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер MongoDB", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер MongoDB успешно запущен")

	// Шаг 3: Запускаем контейнер с приложением
	projectRoot := path.GetProjectRoot()

	// Переопределяем хост MongoDB для подключения к контейнеру из testcontainers
	appEnv := map[string]string{
		testcontainers.MongoHostKey:       generatedMongo.Config().ContainerName,
		testcontainers.MongoPortKey:       testcontainers.MongoPort,
		testcontainers.MongoDatabaseKey:   mongoDatabase,
		testcontainers.MongoUsernameKey:   mongoUsername,
		testcontainers.MongoPasswordKey:   mongoPassword,
		testcontainers.MongoAuthDBKey:     os.Getenv(testcontainers.MongoAuthDBKey),
		testcontainers.MongoMigrationsDir: "/app/migrations",
		grpcPortKey:                       grpcPort,
		testcontainers.MongoLoggerLevel:   loggerLevelValue,
		grpcHostKey:                       grpcHost,
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
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения успешно запущен")

	logger.Info(ctx, "🎉 Тестовое окружение готово")
	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		App:     appContainer,
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
