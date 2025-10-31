package services

import (
	"reception-api/fhir"
	"reception-api/models"
)

type PractitionerService struct {
	fhirClient *fhir.FHIRClient
}

func NewPractitionerService(fhirClient *fhir.FHIRClient) *PractitionerService {
	return &PractitionerService{
		fhirClient: fhirClient,
	}
}

func (s *PractitionerService) GetAllPractitioners() ([]models.PractitionerDTO, error) {
	return s.fhirClient.GetPractitioners()
}
