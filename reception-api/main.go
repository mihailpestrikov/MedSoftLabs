package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reception-api/config"
	"reception-api/database"
	"reception-api/handlers"
	"reception-api/middleware"
	"reception-api/router"
	"reception-api/services"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	repo := database.New(db)

	accessExpiry, _ := time.ParseDuration(cfg.JWTAccessExpiry)
	refreshExpiry, _ := time.ParseDuration(cfg.JWTRefreshExpiry)
	jwtService := middleware.New(cfg.JWTSecret, accessExpiry, refreshExpiry)

	authService := services.NewAuthService(repo, jwtService)
	patientService := services.NewPatientService(repo)

	authHandler := handlers.NewAuthHandler(authService)
	patientHandler := handlers.NewPatientHandler(patientService)

	r := router.Setup(authHandler, patientHandler, jwtService)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		log.Printf("Starting HTTPS server on port %s", cfg.ServerPort)
		if err := srv.ListenAndServeTLS(cfg.TLSCertPath, cfg.TLSKeyPath); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
