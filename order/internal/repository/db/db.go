package db

import (
	"context"
	"database/sql"
)

type DB interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type sqlDB struct {
	*sql.DB
}

func NewSQLDB(db *sql.DB) DB {
	return &sqlDB{db}
}
