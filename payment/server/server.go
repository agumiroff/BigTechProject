package internal

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	payment "github.com/agumiroff/BigTechProject/payment/v1/service"
	payServV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

const (
	port = 50052
)

func StartGRPCServer() (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("Failed to start listener %v\n", err)
		return nil, nil, err
	}

	s := grpc.NewServer()
	service := payment.NewService()
	payServV1.RegisterPaymentServiceServer(s, service)
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
