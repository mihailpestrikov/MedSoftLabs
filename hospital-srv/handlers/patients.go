package handlers

import (
	"hospital-srv/models"
	"hospital-srv/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	service *services.PatientService
}

func New(service *services.PatientService) *PatientHandler {
	return &PatientHandler{service: service}
}

func (h *PatientHandler) GetAllPatients(c *gin.Context) {
	patients, err := h.service.GetAllPatients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, patients)
}

func (h *PatientHandler) GetPatient(c *gin.Context) {
	id := c.Param("id")
	patient, err := h.service.GetPatientByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}
	c.JSON(http.StatusOK, patient)
}

func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var patient models.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.CreatePatient(patient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	patient.ID = id
	c.JSON(http.StatusCreated, patient)
}

func (h *PatientHandler) DeletePatient(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeletePatient(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "patient deleted successfully"})
}

func (h *PatientHandler) BatchDeletePatients(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids array cannot be empty"})
		return
	}

	if err := h.service.DeletePatients(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "patients deleted successfully", "count": len(req.IDs)})
}
