package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reception-api/fhir"
	"reception-api/websocket"

	"github.com/gin-gonic/gin"
)

type FHIRNotificationHandler struct {
	hub *websocket.Hub
}

func NewFHIRNotificationHandler(hub *websocket.Hub) *FHIRNotificationHandler {
	return &FHIRNotificationHandler{
		hub: hub,
	}
}

func (h *FHIRNotificationHandler) HandleEncounterNotification(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read notification body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	var notification map[string]interface{}
	if err := json.Unmarshal(body, &notification); err != nil {
		log.Printf("Failed to unmarshal notification: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification format"})
		return
	}

	eventType, ok := notification["type"].(string)
	if !ok {
		log.Printf("Notification missing type field")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing notification type"})
		return
	}

	encounterData, ok := notification["data"]
	if !ok {
		log.Printf("Notification missing data field")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing notification data"})
		return
	}

	log.Printf("Received FHIR notification: type=%s", eventType)

	dto, err := fhir.MapFHIRToEncounterDTO(encounterData)
	if err != nil {
		log.Printf("Error mapping FHIR to DTO: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to map FHIR data"})
		return
	}

	switch eventType {
	case "encounter_created":
		h.hub.BroadcastEncounterCreated(dto)
	case "encounter_status_updated":
		h.hub.BroadcastEncounterStatusUpdated(dto)
	default:
		log.Printf("Unknown notification type: %s", eventType)
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
