package services

import (
	"hospital-srv/models"
	"hospital-srv/repository"
)

type PractitionerService struct {
	repo *repository.Repository
}

func NewPractitionerService(repo *repository.Repository) *PractitionerService {
	return &PractitionerService{
		repo: repo,
	}
}

func (s *PractitionerService) GetAllPractitioners() ([]models.Practitioner, error) {
	return s.repo.GetAllPractitioners()
}

func (s *PractitionerService) GetPractitionerByID(id string) (*models.Practitioner, error) {
	return s.repo.GetPractitionerByID(id)
}

func (s *PractitionerService) CreatePractitioner(p models.Practitioner) (string, error) {
	return s.repo.CreatePractitioner(p)
}
