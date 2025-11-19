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

	api "github.com/agumiroff/BigTechProject/payment/v1/internal/api/v1"
	"github.com/agumiroff/BigTechProject/payment/v1/internal/config"
	repository "github.com/agumiroff/BigTechProject/payment/v1/internal/repository/payment"
	service "github.com/agumiroff/BigTechProject/payment/v1/internal/service/payment"
	"github.com/agumiroff/BigTechProject/payment/v1/migrations"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

func StartGRPCServer(ctx context.Context) {
	dbURI, dbName, grpcAdress, migPath := loadEnv()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("❌ Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Ping the database
	if err = client.Ping(ctx, nil); err != nil {
		// Explicitly disconnect before returning to avoid defer issue
		if disconnectErr := client.Disconnect(ctx); disconnectErr != nil {
			log.Printf("❌ Warning: Failed to disconnect from MongoDB: %v", disconnectErr)
		}
		log.Printf("❌ Failed to ping MongoDB: %v", err)
		return
	}
	log.Println("✅ Connected to MongoDB")

	db := client.Database(dbName)

	// Run migrations
	if err = migrations.ApplyMigrations(ctx, db, migPath); err != nil {
		log.Printf("Warning: Setup error: %v", err)
	}

	// Initialize repository and service layers
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	api := api.NewAPI(svc)

	// Set up gRPC server
	lis, err := net.Listen("tcp", grpcAdress)
	if err != nil {
		log.Printf("❌ Failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()
	paymentv1.RegisterPaymentServiceServer(s, api)
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
	log.Println("⏹ Shutting down gRPC server")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}

func loadEnv() (string, string, string, string) {
	if err := config.Load(); err != nil {
		log.Fatalf("❌ failed to load config: %v", err)
	}
	dbURI := config.AppConfig().Mongo.URI()
	if dbURI == "" {
		log.Fatal("❌ failed to get MONGO_URI from environment")
	}

	dbName := config.AppConfig().Mongo.DBName()
	if dbName == "" {
		log.Fatal("❌ failed to get MONGO_INITDB_DATABASE from environment")
	}

	grpcAdress := config.AppConfig().GRPC.Address()
	if grpcAdress == "" {
		log.Fatal("❌ failed to get grpc from environment")
	}

	migPath := config.AppConfig().Mongo.MigrationPath()

	return dbURI, dbName, grpcAdress, migPath
}
