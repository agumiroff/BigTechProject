package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-faster/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	invServer "github.com/agumiroff/BigTechProject/inventory/v1/server"
	ExRepo "github.com/agumiroff/BigTechProject/order/v1/external/repository/order"
	"github.com/agumiroff/BigTechProject/order/v1/internal/api/v1"
	handler "github.com/agumiroff/BigTechProject/order/v1/internal/handler/order"
	InRepo "github.com/agumiroff/BigTechProject/order/v1/internal/repository/order"
	serv "github.com/agumiroff/BigTechProject/order/v1/internal/service/order"
	orderServer "github.com/agumiroff/BigTechProject/order/v1/server"
	payServer "github.com/agumiroff/BigTechProject/payment/v1/server"
)

const (
	invServiceAddress = "localhost:50051"
	payServiceAddress = "localhost:50052"
	port              = 8000
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 109 * time.Second
)

func dialGRPC(address string) (conn *grpc.ClientConn) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}

	return conn
}

func closeConn(name string, conn *grpc.ClientConn) {
	if err := conn.Close(); err != nil {
		log.Printf("failed to close %s connection: %v", name, err)
	}
}

func startGRPCServer(consructor func()) {
	consructor()
}

func main() {
	// --- gRPC серверы: INVENTORY и PAYMENT ---
	go func() {
		startGRPCServer(invServer.StartGRPCServer)
	}()

	go func() {
		startGRPCServer(payServer.StartGRPCServer)
	}()

	// --- gRPC клиенты ---
	invConn := dialGRPC(invServiceAddress)
	defer closeConn("Inventory Service", invConn)

	payConn := dialGRPC(payServiceAddress)
	defer closeConn("Payment Service", payConn)

	inRepo := InRepo.NewRepository()
	exRepo := ExRepo.NewRepository(invConn, payConn)
	service := serv.NewService(inRepo, exRepo)
	api := api.NewAPI(service)
	h := handler.NewHandler(api)

	// --- HTTP сервер ---
	server, err := orderServer.StartHTTPServer(h, readHeaderTimeout, port)
	if err != nil {
		log.Printf("❌ Ошибка запуска HTTP-сервера: %v", err)
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %d\n", port)
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
