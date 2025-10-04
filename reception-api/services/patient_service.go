package services

import (
	"errors"
	"reception-api/database"
	"reception-api/models"
)

type PatientService struct {
	repo *database.Repository
}

func NewPatientService(repo *database.Repository) *PatientService {
	return &PatientService{repo: repo}
}

func (s *PatientService) CreatePatient(patient models.Patient) (*models.Patient, error) {
	if patient.FirstName == "" || patient.LastName == "" || patient.DateOfBirth == "" {
		return nil, errors.New("first_name, last_name, and date_of_birth are required")
	}

	id, err := s.repo.CreatePatient(patient)
	if err != nil {
		return nil, err
	}

	patient.ID = id
	return &patient, nil
}

func (s *PatientService) GetAllPatients() ([]models.Patient, error) {
	patients, err := s.repo.GetAllPatients()
	if err != nil {
		return nil, err
	}

	if patients == nil {
		return []models.Patient{}, nil
	}

	return patients, nil
}

func (s *PatientService) GetPatientByID(id int) (*models.Patient, error) {
	patient, err := s.repo.GetPatientByID(id)
	if err != nil {
		return nil, errors.New("patient not found")
	}

	return patient, nil
}

func (s *PatientService) DeletePatient(id int) error {
	_, err := s.repo.GetPatientByID(id)
	if err != nil {
		return errors.New("patient not found")
	}

	return s.repo.DeletePatient(id)
}
