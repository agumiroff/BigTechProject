package main

import (
	"context"
	"log"
	"os"
	"syscall"

	"github.com/agumiroff/BigTechProject/order/v1/internal/app"
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

	// Запускаем HTTP сервер
	go func() {
		if err := a.Run(ctx); err != nil {
			log.Printf("❌ HTTP server error: %v", err)
		}
	}()

	// Ожидаем сигнала завершения
	<-closer.Done()
	log.Println("✅ Service shutdown complete")
}
