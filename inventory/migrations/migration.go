package migrations

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type migKind int16

const (
	UP migKind = iota
	DOWN
)

type Migration struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	Path      string    `bson:"path"`
	AppliedAt time.Time `bson:"appliedAt"`
	Status    string    `bson:"status"`
}

type MigrationError struct {
	MigrationID string    `bson:"migrationId"`
	Name        string    `bson:"name"`
	Path        string    `bson:"path"`
	Error       string    `bson:"error"`
	OccurredAt  time.Time `bson:"occurredAt"`
}

func getMigrationsList(kind migKind, path string) ([]Migration, error) {
	if path == "" {
		return nil, fmt.Errorf("MIGRATIONS_PATH environment variable not set")
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("migrations directory does not exist: %s", path)
		}
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// сортируем по имени (001_*, 002_* …)
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

	migs := make([]Migration, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}

		filename := e.Name()
		ext := filepath.Ext(filename)
		id := strings.TrimSuffix(filename, ext)
		switch kind {

		case UP:
			if !strings.HasSuffix(filename, ".up.json") {
				continue
			}
		case DOWN:
			if !strings.HasSuffix(filename, ".down.json") {
				continue
			}
		}

		migs = append(migs, Migration{
			ID:   id,
			Name: filename,
			Path: filepath.Join(path, filename),
		})
	}
	return migs, nil
}

func ApplyMigrations(ctx context.Context, db *mongo.Database, migPath string) error {
	var kind migKind
	if len(os.Args) > 1 {
		command := os.Args[1]
		if command == "migrate" {
			if os.Args[2] == "up" {
				log.Printf("running up migrations")
				kind = UP
			}
			if os.Args[2] == "down" {
				log.Printf("running down migrations")
				kind = DOWN
			}
		}
	}
	migrations, err := getMigrationsList(kind, migPath)
	if err != nil {
		return fmt.Errorf("failed to get migrations list: %w", err)
	}

	migColl := db.Collection("migrations")
	errColl := db.Collection("migErrors")

	// множество уже применённых миграций
	applied := make(map[string]struct{})
	cur, err := migColl.Find(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer func() {
		if closeErr := cur.Close(ctx); closeErr != nil {
			err = fmt.Errorf("failed to close cursor: %w, original error: %w", closeErr, err)
		}
	}()

	for cur.Next(ctx) {
		var m Migration
		if err := cur.Decode(&m); err != nil {
			return fmt.Errorf("failed to decode applied migration: %w", err)
		}
		applied[m.ID] = struct{}{}
	}
	if err := cur.Err(); err != nil {
		return fmt.Errorf("cursor error while reading applied migrations: %w", err)
	}

	for _, mig := range migrations {
		if _, ok := applied[mig.ID]; ok {
			log.Printf("ℹ️ Skipping already applied migration: %s", mig.Name)
			continue
		}

		log.Printf("🔄 Applying migration: %s", mig.Name)

		data, err := os.ReadFile(mig.Path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", mig.Path, err)
		}

		var cmd bson.Raw
		if err := bson.UnmarshalExtJSON(data, true, &cmd); err != nil {
			return fmt.Errorf("failed to parse migration %s: %w", mig.Name, err)
		}

		if err := db.RunCommand(ctx, cmd).Err(); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", mig.Name, err)
		}

		// migration timeout
		migCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		start := time.Now()
		runErr := db.RunCommand(migCtx, cmd).Err()
		cancel()

		now := time.Now()
		if runErr != nil {
			log.Printf("❌ Migration %s failed: %v", mig.Name, runErr)
			if _, insertErr := errColl.InsertOne(ctx, MigrationError{
				MigrationID: mig.ID,
				Name:        mig.Name,
				Path:        mig.Path,
				Error:       runErr.Error(),
				OccurredAt:  now,
			}); insertErr != nil {
				log.Printf("❌ Failed to log migration error to database: %v", insertErr)
			}
			return fmt.Errorf("failed to apply migration %s: %w", mig.Name, runErr)
		}

		_, insErr := migColl.InsertOne(ctx, Migration{
			ID:        mig.ID,
			Name:      mig.Name,
			Path:      mig.Path,
			AppliedAt: now.UTC(),
			Status:    "applied",
		})
		if insErr != nil && !mongo.IsDuplicateKeyError(insErr) {
			return fmt.Errorf("failed to record migration %s (took %s): %w", mig.Name, time.Since(start), insErr)
		}

		log.Printf("✅ Successfully applied migration: %s (took %s)", mig.Name, time.Since(start))
	}

	return nil
}
