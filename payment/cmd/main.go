package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/agumiroff/BigTechProject/payment/v1/server"
)

func main() {
	ctx := context.Background()

	// Load .env file from project root
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	dbURI := os.Getenv("MONGO_URI")
	if dbURI == "" {
		log.Fatal("❌ failed to get MONGO_URI from environment")
	}

	dbName := os.Getenv("MONGO_INITDB_DATABASE")
	if dbName == "" {
		log.Fatal("failed to get MONGO_INITDB_DATABASE from environment")
	}

	server.StartGRPCServer(ctx, dbURI, dbName)
}
