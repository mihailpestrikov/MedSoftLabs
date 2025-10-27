package services

import "reception-api/fhir"

type PractitionerService struct {
	fhirClient *fhir.FHIRClient
}

func NewPractitionerService(fhirClient *fhir.FHIRClient) *PractitionerService {
	return &PractitionerService{
		fhirClient: fhirClient,
	}
}

func (s *PractitionerService) GetAllPractitioners() ([]fhir.Practitioner, error) {
	return s.fhirClient.GetPractitioners()
}
