package database

import (
	"database/sql"
	"fmt"
	"hospital-srv/config"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/lib/pq"
)

func Setup(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")

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

	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	sort.Strings(sqlFiles)

	for _, file := range sqlFiles {
		content, err := os.ReadFile(filepath.Join("migrations", file))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		log.Printf("Migration %s applied successfully", file)
	}

	return nil
}
