package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/agumiroff/BigTechProject/inventory/v1/internal/api/inventory/v1"
	repository "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/part"
	service "github.com/agumiroff/BigTechProject/inventory/v1/internal/service/part"
	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func StartGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to start listener %v\n", err)
		return
	}

	s := grpc.NewServer()

	repo := repository.NewRepository()
	service := service.NewService(repo)
	api := api.NewAPI(service)

	inventoryv1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Printf("Starting grpc server on port %v", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to start server %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down grpc server")
	s.GracefulStop()
	log.Println("Server stopped")
}
