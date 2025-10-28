package handlers

import (
	"net/http"
	"reception-api/services"

	"github.com/gin-gonic/gin"
)

type PractitionerHandler struct {
	practitionerService *services.PractitionerService
}

func NewPractitionerHandler(practitionerService *services.PractitionerService) *PractitionerHandler {
	return &PractitionerHandler{practitionerService: practitionerService}
}

func (h *PractitionerHandler) GetAllPractitioners(c *gin.Context) {
	practitioners, err := h.practitionerService.GetAllPractitioners()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, practitioners)
}
