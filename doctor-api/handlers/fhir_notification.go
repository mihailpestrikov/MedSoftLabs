package handlers

import (
	"doctor-api/fhir"
	"doctor-api/websocket"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FHIRNotificationHandler handles FHIR notifications from HIS.
type FHIRNotificationHandler struct {
	hub *websocket.Hub
}

// NewFHIRNotificationHandler creates a new FHIR notification handler.
func NewFHIRNotificationHandler(hub *websocket.Hub) *FHIRNotificationHandler {
	return &FHIRNotificationHandler{hub: hub}
}

// HandleEncounterNotification processes encounter notifications and broadcasts to connected clients.
func (h *FHIRNotificationHandler) HandleEncounterNotification(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading FHIR notification body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	log.Printf("Received FHIR Encounter notification from HIS: %s", string(body))

	var notification map[string]interface{}
	if err := json.Unmarshal(body, &notification); err != nil {
		log.Printf("Error parsing FHIR notification: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	eventType, ok := notification["type"].(string)
	if !ok {
		log.Printf("Missing or invalid notification type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing notification type"})
		return
	}

	encounterData, ok := notification["data"]
	if !ok {
		log.Printf("Missing notification data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing notification data"})
		return
	}

	dto, err := fhir.MapFHIRToEncounterDTO(encounterData)
	if err != nil {
		log.Printf("Error mapping FHIR to DTO: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to map FHIR data"})
		return
	}

	switch eventType {
	case "encounter_created":
		h.hub.Broadcast(websocket.Message{
			Type: "encounter_created",
			Data: dto,
		})
		log.Printf("Broadcasted encounter_created event to Doctor.UI clients")

	case "encounter_status_updated":
		h.hub.Broadcast(websocket.Message{
			Type: "encounter_status_updated",
			Data: dto,
		})
		log.Printf("Broadcasted encounter_status_updated event to Doctor.UI clients")

	default:
		log.Printf("Unknown notification type: %s", eventType)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification received"})
}
