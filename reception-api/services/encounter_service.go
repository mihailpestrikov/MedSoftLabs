package services

import (
	"errors"
	"fmt"
	"reception-api/database"
	"reception-api/fhir"
	"reception-api/models"
	"time"
)

var validEncounterStatuses = map[string]bool{
	"planned":     true,
	"arrived":     true,
	"in-progress": true,
	"completed":   true,
	"cancelled":   true,
}

type EncounterService struct {
	repo       *database.Repository
	fhirClient *fhir.FHIRClient
}

func NewEncounterService(repo *database.Repository, fhirClient *fhir.FHIRClient) *EncounterService {
	return &EncounterService{
		repo:       repo,
		fhirClient: fhirClient,
	}
}

type CreateEncounterRequest struct {
	PatientID      int       `json:"patient_id" binding:"required"`
	PractitionerID string    `json:"practitioner_id" binding:"required"`
	StartTime      time.Time `json:"start_time" binding:"required"`
}

func (s *EncounterService) CreateEncounter(req CreateEncounterRequest) (string, error) {
	patient, err := s.repo.GetPatientByID(req.PatientID)
	if err != nil {
		return "", errors.New("patient not found")
	}

	if patient.HISPatientID == nil || *patient.HISPatientID == "" {
		return "", errors.New("patient does not have HIS Patient ID yet")
	}

	encounterID, err := s.fhirClient.CreateEncounter(*patient.HISPatientID, req.PractitionerID, req.StartTime)
	if err != nil {
		return "", err
	}

	return encounterID, nil
}

func (s *EncounterService) UpdateEncounterStatus(encounterID string, status string) error {
	if !validEncounterStatuses[status] {
		return fmt.Errorf("invalid encounter status: %s (valid statuses: planned, arrived, in-progress, completed, cancelled)", status)
	}
	return s.fhirClient.UpdateEncounterStatus(encounterID, status)
}

func (s *EncounterService) GetAllEncounters() ([]models.EncounterDTO, error) {
	return s.fhirClient.GetEncounters()
}
