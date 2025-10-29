package handlers

import (
	"doctor-api/fhir"
	"doctor-api/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EncounterHandler struct {
	fhirClient *fhir.FHIRClient
	hub        *websocket.Hub
}

func NewEncounterHandler(fhirClient *fhir.FHIRClient, hub *websocket.Hub) *EncounterHandler {
	return &EncounterHandler{
		fhirClient: fhirClient,
		hub:        hub,
	}
}

func (h *EncounterHandler) GetEncountersByPractitioner(c *gin.Context) {
	practitionerID := c.Param("practitioner_id")
	if practitionerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "practitioner_id is required"})
		return
	}

	encounters, err := h.fhirClient.GetEncountersByPractitioner(practitionerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, encounters)
}

func (h *EncounterHandler) UpdateEncounterStatus(c *gin.Context) {
	encounterID := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.fhirClient.UpdateEncounterStatus(encounterID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}
