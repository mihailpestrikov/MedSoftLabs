package main

import (
	"doctor-api/config"
	"doctor-api/fhir"
	"doctor-api/handlers"
	"doctor-api/router"
	"doctor-api/websocket"
	"fmt"
	"log"
)

func main() {
	cfg := config.Load()

	fhirClient, err := fhir.NewFHIRClient(cfg.HISHTTPAddress, cfg.TLSCertPath)
	if err != nil {
		log.Fatalf("Failed to create FHIR client: %v", err)
	}

	hub := websocket.NewHub()
	go hub.Run()

	encounterHandler := handlers.NewEncounterHandler(fhirClient, hub)
	practitionerHandler := handlers.NewPractitionerHandler(fhirClient)
	fhirNotificationHandler := handlers.NewFHIRNotificationHandler(hub)

	r := router.Setup(encounterHandler, practitionerHandler, fhirNotificationHandler, hub)

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Starting Doctor API server on %s", serverAddr)
	if err := r.RunTLS(serverAddr, cfg.TLSCertPath, cfg.TLSKeyPath); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
