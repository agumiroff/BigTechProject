package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-faster/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	exRepo "github.com/agumiroff/BigTechProject/order/v1/external/repository/order"
	"github.com/agumiroff/BigTechProject/order/v1/internal/api/v1"
	"github.com/agumiroff/BigTechProject/order/v1/internal/config"
	"github.com/agumiroff/BigTechProject/order/v1/internal/db"
	handler "github.com/agumiroff/BigTechProject/order/v1/internal/handler/order"
	"github.com/agumiroff/BigTechProject/order/v1/internal/migrator"
	inRepo "github.com/agumiroff/BigTechProject/order/v1/internal/repository/order"
	serv "github.com/agumiroff/BigTechProject/order/v1/internal/service/order"
	orderServer "github.com/agumiroff/BigTechProject/order/v1/server"
)

func dialGRPC(address string) (conn *grpc.ClientConn) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return conn
	}

	return conn
}

func closeConn(name string, conn *grpc.ClientConn) {
	if err := conn.Close(); err != nil {
		log.Printf("failed to close %s connection: %v", name, err)
	}
}

func main() {
	err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	httpConfig := config.AppConfig().HTTPConfig
	if httpConfig == nil {
		log.Fatalf("HTTP configuration is nil after loading")
	}

	// envs loading
	invServiceAddress := httpConfig.InventoryAddr()
	payServiceAddress := httpConfig.PaymentAddr()
	readHeaderTimeout := httpConfig.Timeout()
	address := httpConfig.Address()
	shutdownTimeout := httpConfig.Timeout()
	migrationsDir := config.AppConfig().PostgressConfig.MigPath()
	dbURI := config.AppConfig().PostgressConfig.DSN()

	// --- gRPC clients ---
	invConn := dialGRPC(invServiceAddress)
	if invConn != nil {
		defer closeConn("Inventory Service", invConn)
	}

	payConn := dialGRPC(payServiceAddress)
	if payConn != nil {
		defer closeConn("Payment Service", payConn)
	}

	// Create repositories
	db, err := db.ConnectDB(dbURI)
	if err != nil {
		log.Printf("failed to create database connection %v", err)
		return
	}
	m := migrator.NewMigrator(db, migrationsDir)
	if err := m.RunMigrations(); err != nil {
		log.Printf("Failed to run migrations: %v", err)
	}

	inRepo := inRepo.NewRepository(db)
	exRepo := exRepo.NewRepository(invConn, payConn)
	service := serv.NewService(inRepo, exRepo)
	api := api.NewAPI(service)
	h := handler.NewHandler(api)

	// --- HTTP сервер ---
	server, err := orderServer.StartHTTPServer(h, readHeaderTimeout, address)
	if err != nil {
		log.Printf("❌ Ошибка запуска HTTP-сервера: %v", err)
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", address)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка работы HTTP-сервера: %v", err)
		}
	}()

	// --- Graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Завершение работы серверов...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Printf("❌ Ошибка при остановке HTTP-сервера: %v", err)
	}

	log.Println("✅ HTTP-сервер остановлен корректно")
}
