package handlers

import (
	"net/http"
	"reception-api/services"

	"github.com/gin-gonic/gin"
)

type EncounterHandler struct {
	encounterService *services.EncounterService
}

func NewEncounterHandler(encounterService *services.EncounterService) *EncounterHandler {
	return &EncounterHandler{encounterService: encounterService}
}

func (h *EncounterHandler) CreateEncounter(c *gin.Context) {
	var req services.CreateEncounterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encounterID, err := h.encounterService.CreateEncounter(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "encounter created successfully",
		"encounter_id": encounterID,
	})
}
