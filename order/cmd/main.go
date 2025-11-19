package main

import (
	"context"
	"os"
	"syscall"

	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/order/v1/internal/app"
	"github.com/agumiroff/BigTechProject/platform/closer"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
)

func main() {
	ctx := context.Background()

	// Конфигурируем closer для graceful shutdown
	closer.Configure(syscall.SIGTERM, os.Interrupt)

	// Создаём и инициализируем приложение
	a, err := app.New(ctx)
	if err != nil {
		logger.Fatal(ctx, "❌ Failed to initialize app", zap.Error(err))
	}

	// Запускаем HTTP сервер
	go func() {
		if err := a.Run(ctx); err != nil {
			logger.Error(ctx, "❌ HTTP server error", zap.Error(err))
		}
	}()

	// Ожидаем сигнала завершения
	<-closer.Done()
	logger.Info(ctx, "✅ Service shutdown complete")
}
