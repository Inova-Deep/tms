package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

func Init(dbPath string) (*sql.DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Connected to database at %s", dbPath)

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {
		version := filepath.Base(file[:len(file)-4])

		var exists bool
		err := db.QueryRow(
			"SELECT 1 FROM schema_migrations WHERE version = ?",
			version,
		).Scan(&exists)
		if err == nil {
			log.Printf("Skipping migration (already applied): %s", file)
			continue
		}
		if err != sql.ErrNoRows {
			return fmt.Errorf("failed to check migration status %s: %w", file, err)
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		if _, err := db.Exec(
			"INSERT OR IGNORE INTO schema_migrations (version) VALUES (?)",
			version,
		); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}

		log.Printf("Applied migration: %s", file)
	}
	return nil
}
