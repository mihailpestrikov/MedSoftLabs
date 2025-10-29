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
		"id": encounterID,
	})
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

	if err := h.encounterService.UpdateEncounterStatus(encounterID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

func (h *EncounterHandler) GetAllEncounters(c *gin.Context) {
	encounters, err := h.encounterService.GetAllEncounters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, encounters)
}
