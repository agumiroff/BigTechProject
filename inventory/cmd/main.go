package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	invserver "github.com/agumiroff/BigTechProject/inventory/v1/server"
)

func main() {
	s, lis, err := invserver.StartGRPCServer()
	if err != nil {
		log.Fatalf("inventory server error: %v", err)
	}

	defer func() {
		if err := lis.Close(); err != nil {
			log.Printf("failed to close listener: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down GRPC server")
	s.GracefulStop()
	log.Println("Server stopped")
}
