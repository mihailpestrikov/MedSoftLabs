package services

import (
	"hospital-srv/models"
	"hospital-srv/repository"
	"hospital-srv/websocket"
)

type PatientService struct {
	repo *repository.Repository
	hub  *websocket.Hub
}

func New(repo *repository.Repository, hub *websocket.Hub) *PatientService {
	return &PatientService{
		repo: repo,
		hub:  hub,
	}
}

func (s *PatientService) CreatePatient(patient models.Patient) (string, error) {
	id, err := s.repo.CreatePatient(patient)
	if err != nil {
		return "", err
	}

	createdPatient, err := s.repo.GetPatientByID(id)
	if err != nil {
		return id, nil
	}

	s.hub.BroadcastPatientCreated(createdPatient)

	return id, nil
}

func (s *PatientService) GetAllPatients() ([]models.Patient, error) {
	return s.repo.GetAllPatients()
}

func (s *PatientService) GetPatientByID(id string) (*models.Patient, error) {
	return s.repo.GetPatientByID(id)
}

func (s *PatientService) DeletePatient(id string) error {
	if err := s.repo.DeletePatient(id); err != nil {
		return err
	}

	s.hub.BroadcastPatientDeleted(id)

	return nil
}

func (s *PatientService) DeletePatients(ids []string) error {
	if err := s.repo.DeletePatients(ids); err != nil {
		return err
	}

	for _, id := range ids {
		s.hub.BroadcastPatientDeleted(id)
	}

	return nil
}
