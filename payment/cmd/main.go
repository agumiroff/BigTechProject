package main

import (
	"context"
	"log"
	"os"
	"syscall"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/app"
	"github.com/agumiroff/BigTechProject/platform/closer"
)

func main() {
	ctx := context.Background()

	// Конфигурируем closer для graceful shutdown
	closer.Configure(syscall.SIGTERM, os.Interrupt)

	// Создаём и инициализируем приложение
	a, err := app.New(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to initialize app: %v", err)
	}

	// Запускаем gRPC сервер
	go func() {
		if err := a.Run(ctx); err != nil {
			log.Printf("❌ gRPC server error: %v", err)
		}
	}()

	// Ожидаем сигнала завершения
	<-closer.Done()
	log.Println("✅ Service shutdown complete")
}
