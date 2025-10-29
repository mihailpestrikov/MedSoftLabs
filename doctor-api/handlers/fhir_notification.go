package handlers

import (
	"doctor-api/websocket"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FHIRNotificationHandler struct {
	hub *websocket.Hub
}

func NewFHIRNotificationHandler(hub *websocket.Hub) *FHIRNotificationHandler {
	return &FHIRNotificationHandler{hub: hub}
}

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

	eventType, _ := notification["type"].(string)
	encounterData, _ := notification["data"]

	switch eventType {
	case "encounter_created":
		h.hub.Broadcast(websocket.Message{
			Type: "encounter_created",
			Data: encounterData,
		})
		log.Printf("Broadcasted encounter_created event to Doctor.UI clients")

	case "encounter_status_updated":
		h.hub.Broadcast(websocket.Message{
			Type: "encounter_status_updated",
			Data: encounterData,
		})
		log.Printf("Broadcasted encounter_status_updated event to Doctor.UI clients")

	default:
		log.Printf("Unknown notification type: %s", eventType)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification received"})
}
