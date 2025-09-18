package database

import (
	"context"
	"database/sql"
)

type DB interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type SQLDB struct {
	*sql.DB
}

func NewSQLDB(db *sql.DB) DB {
	return &SQLDB{db}
}
