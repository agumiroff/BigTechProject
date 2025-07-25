package server

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/part"
	invServiceV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

const (
	port = 50051
)

func StartGRPCServer() (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("failed to start listener: %v", err)
		return nil, nil, err
	}

	s := grpc.NewServer()

	repo := repository.NewRepository()
	invServiceV1.RegisterInventoryServiceServer(s, repo)
	reflection.Register(s)

	go func() {
		log.Printf("Starting grpc-server on port %d\n", port)
		err := s.Serve(lis)
		if err != nil {
			log.Printf("Failed to start grpc-server: %v\n", err)
			return
		}
	}()

	return s, lis, nil
}
