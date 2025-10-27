package main

import (
	"context"
	"errors"
	"hospital-srv/config"
	"hospital-srv/database"
	"hospital-srv/fhir"
	"hospital-srv/handlers"
	"hospital-srv/hl7"
	"hospital-srv/repository"
	"hospital-srv/router"
	"hospital-srv/services"
	"hospital-srv/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	repo := repository.New(db)
	patientService := services.New(repo, hub)
	practitionerService := services.NewPractitionerService(repo)
	encounterService := services.NewEncounterService(repo, hub)

	patientHandler := handlers.New(patientService)
	fhirServer := fhir.NewFHIRServer(practitionerService, encounterService)

	r := router.Setup(patientHandler, hub, fhirServer)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	hl7Handler := hl7.NewHL7Handler(patientService)
	mllpListener, err := hl7.NewMLLPListener(cfg.MLLPPort, cfg.TLSCertPath, cfg.TLSKeyPath, hl7Handler.HandleMessage)
	if err != nil {
		log.Fatalf("Failed to start MLLP listener: %v", err)
	}
	defer mllpListener.Close()

	go mllpListener.Start()

	go func() {
		log.Printf("HTTP server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
