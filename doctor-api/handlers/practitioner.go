package handlers

import (
	"doctor-api/fhir"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PractitionerHandler struct {
	fhirClient *fhir.FHIRClient
}

func NewPractitionerHandler(fhirClient *fhir.FHIRClient) *PractitionerHandler {
	return &PractitionerHandler{fhirClient: fhirClient}
}

func (h *PractitionerHandler) GetAllPractitioners(c *gin.Context) {
	practitioners, err := h.fhirClient.GetPractitioners()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, practitioners)
}
