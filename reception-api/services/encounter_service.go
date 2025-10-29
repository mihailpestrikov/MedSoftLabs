package services

import (
	"errors"
	"reception-api/database"
	"reception-api/fhir"
	"time"
)

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
	return s.fhirClient.UpdateEncounterStatus(encounterID, status)
}

func (s *EncounterService) GetAllEncounters() ([]map[string]interface{}, error) {
	return s.fhirClient.GetEncounters()
}
