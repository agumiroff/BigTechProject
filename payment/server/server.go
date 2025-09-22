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

	api "github.com/agumiroff/BigTechProject/payment/v1/internal/api/v1"
	repo "github.com/agumiroff/BigTechProject/payment/v1/internal/repository/payment"
	service "github.com/agumiroff/BigTechProject/payment/v1/internal/service/payment"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

func StartGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to start listener %v\n", err)
		return
	}

	s := grpc.NewServer()

	repo := repo.NewRepository()
	service := service.NewService(repo)
	api := api.NewAPI(service)

	paymentv1.RegisterPaymentServiceServer(s, api)

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
