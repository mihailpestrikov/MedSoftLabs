package services

import (
	"hospital-srv/database"
	"hospital-srv/models"
)

type PatientService struct {
	repo *database.Repository
}

func New(repo *database.Repository) *PatientService {
	return &PatientService{repo: repo}
}

func (s *PatientService) CreatePatient(patient models.Patient) (string, error) {
	return s.repo.CreatePatient(patient)
}

func (s *PatientService) GetAllPatients() ([]models.Patient, error) {
	return s.repo.GetAllPatients()
}

func (s *PatientService) GetPatientByID(id string) (*models.Patient, error) {
	return s.repo.GetPatientByID(id)
}

func (s *PatientService) DeletePatient(id string) error {
	return s.repo.DeletePatient(id)
}
