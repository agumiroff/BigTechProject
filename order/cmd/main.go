package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-faster/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	invServer "github.com/agumiroff/BigTechProject/inventory/v1/server"
	"github.com/agumiroff/BigTechProject/order/v1/handler"
	orderServer "github.com/agumiroff/BigTechProject/order/v1/server"
	payServer "github.com/agumiroff/BigTechProject/payment/v1/server"
	InvV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
	PayV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

const (
	invServiceAddress = "localhost:50051"
	payServiceAddress = "localhost:50052"
	address           = "localhost"
	port              = 800
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
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

func startGRPCServer(name string, constructor func() (*grpc.Server, net.Listener, error)) {
	server, lis, err := constructor()
	if err != nil {
		log.Fatalf("❌ Ошибка запуска %s GRPC сервера: %v", name, err)
	}
	defer func(lis net.Listener) {
		err := lis.Close()
		if err != nil {
			log.Printf("Failed to close %s server listener %d\n", name, err)
		}
	}(lis)
	log.Printf("✅ %s gRPC сервер запущен", name)
	if err := server.Serve(lis); err != nil {
		log.Printf("❌ %s server serve error: %v", name, err)
	}
}

func main() {
	// --- gRPC серверы: INVENTORY и PAYMENT ---
	go func() {
		startGRPCServer("inventory", invServer.StartGRPCServer)
	}()

	go func() {
		startGRPCServer("payment", payServer.StartGRPCServer)
	}()

	// --- gRPC клиенты ---
	invConn := dialGRPC(invServiceAddress)
	defer closeConn("Inventory Service", invConn)

	payConn := dialGRPC(payServiceAddress)
	defer closeConn("Payment Service", payConn)

	invClient := InvV1.NewInventoryServiceClient(invConn)
	payClient := PayV1.NewPaymentServiceClient(payConn)
	h := handler.NewOrderHandler(invClient, payClient)

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
