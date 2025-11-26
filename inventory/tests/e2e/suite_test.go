package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
)

const testsTimeout = 5 * time.Minute

var (
	env         *TestEnvironment
	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UFO Service Integration Test Suite")
}

var _ = BeforeSuite(func() {
	err := logger.Init(loggerLevelValue, true)
	if err != nil {
		panic(fmt.Sprintf("не удалось инициализировать логгер: %v", err))
	}

	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	// Загружаем .env файл и устанавливаем переменные в окружение
	envVars, err := godotenv.Read(filepath.Join(getProjectRoot(), "deploy", "compose", "inventory", ".env"))
	if err != nil {
		logger.Fatal(suiteCtx, "Не удалось загрузить .env файл", zap.Error(err))
	}

	// Устанавливаем переменные в окружение процесса
	for key, value := range envVars {
		_ = os.Setenv(key, value)
	}

	logger.Info(suiteCtx, "Запуск тестового окружения...")
	env = setupTestEnvironment(suiteCtx)
})

var _ = AfterSuite(func() {
	logger.Info(context.Background(), "Завершение набора тестов")
	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}
	suiteCancel()
})

func getProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("не удалось получить рабочую директорию: " + err.Error())
	}

	for {
		_, err = os.Stat(filepath.Join(dir, "go.work"))
		if err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			panic("не удалось найти корень проекта (go.work)")
		}

		dir = parent
	}
}
