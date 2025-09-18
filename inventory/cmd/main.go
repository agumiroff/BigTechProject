package main

import (
	"context"
	"log"
	"os"

	"github.com/agumiroff/BigTechProject/inventory/v1/server"
	"github.com/joho/godotenv"
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

	// Get database name from environment
	dbName := os.Getenv("MONGO_INITDB_DATABASE")
	if dbName == "" {
		log.Fatal("failed to get MONGO_INITDB_DATABASE from environment")
	}

	server.StartGRPCServer(ctx, dbURI, dbName)
}
