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
	migration, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration: %w", err)
	}

	if _, err := db.Exec(string(migration)); err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
