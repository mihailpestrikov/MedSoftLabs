package services

import (
	"hospital-srv/models"
	"hospital-srv/repository"
	"hospital-srv/websocket"
)

type EncounterService struct {
	repo *repository.Repository
	hub  *websocket.Hub
}

func NewEncounterService(repo *repository.Repository, hub *websocket.Hub) *EncounterService {
	return &EncounterService{
		repo: repo,
		hub:  hub,
	}
}

func (s *EncounterService) CreateEncounter(encounter models.Encounter) (string, error) {
	id, err := s.repo.CreateEncounter(encounter)
	if err != nil {
		return "", err
	}

	createdEncounter, err := s.repo.GetEncounterByID(id)
	if err != nil {
		return id, nil
	}

	s.hub.BroadcastEncounterCreated(createdEncounter)

	return id, nil
}

func (s *EncounterService) GetAllEncounters() ([]models.EncounterWithDetails, error) {
	return s.repo.GetAllEncounters()
}

func (s *EncounterService) GetEncounterByID(id string) (*models.EncounterWithDetails, error) {
	return s.repo.GetEncounterByID(id)
}
