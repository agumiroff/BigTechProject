package server

import (
	"context"
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
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/config"
	repository "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/part"
	service "github.com/agumiroff/BigTechProject/inventory/v1/internal/service/part"
	"github.com/agumiroff/BigTechProject/inventory/v1/migrations"
	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

func StartGRPCServer(ctx context.Context) {
	dbURI, dbName, grpcAdress, migPath := loadEnv()
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("❌ Failed to connect to MongoDB: %v", err)
		return
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("❌ Warning: Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Ping the database
	if err = client.Ping(ctx, nil); err != nil {
		log.Printf("❌ Failed to ping MongoDB: %v", err)
		// Return early instead of using log.Fatalf, allowing defer to run
		return
	}
	log.Println("✅ Connected to MongoDB")

	db := client.Database(dbName)

	// Run migrations
	if err = migrations.ApplyMigrations(ctx, db, migPath); err != nil {
		log.Printf("Warning: Setup error: %v", err)
	}
	log.Printf("Migrations applies")

	// Initialize repository and service layers
	repo := repository.NewRepository(ctx, db)
	svc := service.NewService(repo)
	api := api.NewAPI(svc)

	// Set up gRPC server
	lis, err := net.Listen("tcp", grpcAdress)
	if err != nil {
		log.Printf("❌ Failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()
	inventoryv1.RegisterInventoryServiceServer(s, api)
	reflection.Register(s)

	// Handle shutdown gracefully
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// Channel for server errors
	errCh := make(chan error, 1)

	go func() {
		log.Printf("🚀 Starting gRPC server on port %v", grpcAdress)
		if err := s.Serve(lis); err != nil {
			log.Printf("❌ Failed to serve: %v", err)
			errCh <- err
		}
	}()

	// Handle either quit signal or server error
	select {
	case <-quit:
		// Normal shutdown, continue below
	case serverErr := <-errCh:
		log.Printf("❌ Server error: %v", serverErr)
		s.Stop() // Force stop in case of error
		return
	}

	// Moved selection logic above
	log.Println("⏹ Shutting down gRPC server")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}

func loadEnv() (string, string, string, string) {
	config.Load()
	dbURI := config.AppConfig().MongoConfig.URI()
	if dbURI == "" {
		log.Fatal("❌ failed to get MONGO_URI from environment")
	}

	dbName := config.AppConfig().MongoConfig.DBName()
	if dbName == "" {
		log.Fatal("❌ failed to get MONGO_INITDB_DATABASE from environment")
	}

	grpcAdress := config.AppConfig().GRPCConfig.Address()
	if grpcAdress == "" {
		log.Fatal("❌ failed to get gRPC address from environment")
	}

	migPath := config.AppConfig().MongoConfig.MigrationPath()

	return dbURI, dbName, grpcAdress, migPath
}
