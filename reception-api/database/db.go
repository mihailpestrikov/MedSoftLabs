package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reception-api/config"

	_ "github.com/lib/pq"
)

func Setup(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Printf("Connected to PostgreSQL: %s:%s/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)

	if err := RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	files, err := os.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && len(file.Name()) > 4 && file.Name()[len(file.Name())-4:] == ".sql" {
			migrationFiles = append(migrationFiles, "migrations/"+file.Name())
		}
	}

	if len(migrationFiles) == 0 {
		log.Println("No migrations found")
		return nil
	}

	for _, migrationFile := range migrationFiles {
		migration, err := os.ReadFile(migrationFile)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", migrationFile, err)
		}

		if _, err := db.Exec(string(migration)); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", migrationFile, err)
		}

		log.Printf("Applied migration: %s", migrationFile)
	}

	log.Println("All migrations applied successfully")
	return nil
}
