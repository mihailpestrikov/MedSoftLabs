package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reception-api/config"
	"reception-api/database"
	"reception-api/fhir"
	"reception-api/handlers"
	"reception-api/hl7"
	"reception-api/middleware"
	"reception-api/router"
	"reception-api/services"
	"reception-api/websocket"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	db, err := database.Setup(cfg)
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}
	defer db.Close()

	hub := websocket.NewHub()
	go hub.Run()

	repo := database.New(db)

	accessExpiry, _ := time.ParseDuration(cfg.JWTAccessExpiry)
	refreshExpiry, _ := time.ParseDuration(cfg.JWTRefreshExpiry)
	jwtService := middleware.New(cfg.JWTSecret, accessExpiry, refreshExpiry)

	mllpClient, err := hl7.NewMLLPClient(cfg.HISAddress, cfg.TLSCertPath)
	if err != nil {
		log.Fatalf("Failed to create MLLP client: %v", err)
	}

	fhirClient, err := fhir.NewFHIRClient("https://"+cfg.HISHTTPAddress, cfg.TLSCertPath)
	if err != nil {
		log.Fatalf("Failed to create FHIR client: %v", err)
	}

	authService := services.NewAuthService(repo, jwtService)
	patientService := services.NewPatientService(repo, hub, mllpClient)
	encounterService := services.NewEncounterService(repo, fhirClient)
	practitionerService := services.NewPractitionerService(fhirClient)

	authHandler := handlers.NewAuthHandler(authService)
	patientHandler := handlers.NewPatientHandler(patientService)
	encounterHandler := handlers.NewEncounterHandler(encounterService)
	practitionerHandler := handlers.NewPractitionerHandler(practitionerService)
	fhirNotificationHandler := handlers.NewFHIRNotificationHandler(hub)

	r := router.Setup(authHandler, patientHandler, encounterHandler, practitionerHandler, fhirNotificationHandler, jwtService, hub)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		log.Printf("Starting HTTPS server on port %s", cfg.ServerPort)
		if err := srv.ListenAndServeTLS(cfg.TLSCertPath, cfg.TLSKeyPath); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
