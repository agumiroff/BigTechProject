package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // Register pgx driver
)

func ConnectDB(uri string) (*sql.DB, error) {
	ctx := context.Background()
	if uri == "" {
		return nil, fmt.Errorf("env variable POSTGRESS_URI is empty")
	}

	db, err := sql.Open("pgx", uri)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
