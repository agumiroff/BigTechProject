package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/agumiroff/BigTechProject/inventory/v1/internal/api/inventory/v1"
	repository "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/part"
	service "github.com/agumiroff/BigTechProject/inventory/v1/internal/service/part"
	"github.com/agumiroff/BigTechProject/inventory/v1/migrations"
	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort = 50051
)

func StartGRPCServer(ctx context.Context, dbURI, dbName string) {
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Ping the database
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ Failed to ping MongoDB: %v", err)
	}
	log.Println("✅ Connected to MongoDB")

	db := client.Database(dbName)

	// Run migrations
	if err = migrations.ApplyMigrations(ctx, db); err != nil {
		log.Printf("Warning: Setup error: %v", err)
	}

	// Initialize repository and service layers
	repo := repository.NewRepository(ctx, db)
	svc := service.NewService(repo)
	api := api.NewAPI(svc)

	// Set up gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("❌ Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	inventoryv1.RegisterInventoryServiceServer(s, api)
	reflection.Register(s)

	// Handle shutdown gracefully
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Printf("🚀 Starting gRPC server on port %v", grpcPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("❌ Failed to serve: %v", err)
		}
	}()

	<-quit
	log.Println("⏹ Shutting down gRPC server")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
