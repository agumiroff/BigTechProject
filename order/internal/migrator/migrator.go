package migrator

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

type migrator struct {
	db  *sql.DB
	dir string
}

func NewMigrator(db *sql.DB, dir string) *migrator {
	return &migrator{
		db:  db,
		dir: dir,
	}
}

func (m *migrator) Up() error {
	err := goose.Up(m.db, m.dir)
	if err != nil {
		return err
	}
	return nil
}

func (m *migrator) RunMigrations() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect %w", err)
	}

	if err := goose.Up(m.db, m.dir); err != nil {
		return fmt.Errorf("failed to run up migrations %w", err)
	}

	return nil
}
