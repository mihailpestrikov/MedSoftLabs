package services

import (
	"errors"
	"log"
	"reception-api/database"
	"reception-api/hl7"
	"reception-api/models"
	"reception-api/websocket"
	"strings"
)

type PatientService struct {
	repo       *database.Repository
	hub        *websocket.Hub
	mllpClient *hl7.MLLPClient
}

func NewPatientService(repo *database.Repository, hub *websocket.Hub, mllpClient *hl7.MLLPClient) *PatientService {
	return &PatientService{
		repo:       repo,
		hub:        hub,
		mllpClient: mllpClient,
	}
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

	s.hub.BroadcastPatientCreated(&patient)

	go s.sendToHIS(&patient)

	return &patient, nil
}

func (s *PatientService) sendToHIS(patient *models.Patient) {
	messageID, hl7Message := hl7.GenerateADTA04(patient)

	log.Printf("Sending patient %d to HIS via HL7 (MessageID: %s)", patient.ID, messageID)
	log.Printf("HL7 ADT^A04: %s", strings.ReplaceAll(string(hl7Message), "\r", "|"))

	ack, err := s.mllpClient.SendMessage(hl7Message)
	if err != nil {
		log.Printf("Failed to send HL7 message: %v", err)
		return
	}

	log.Printf("Received ACK: %s", strings.ReplaceAll(string(ack), "\r", "|"))

	originalMsgID, hisPatientID, err := hl7.ParseACK(ack)
	if err != nil {
		log.Printf("Failed to parse ACK: %v", err)
		return
	}

	if originalMsgID != messageID {
		log.Printf("ACK MessageID mismatch: expected %s, got %s", messageID, originalMsgID)
		return
	}

	log.Printf("Received HIS Patient ID: %s for local patient %d", hisPatientID, patient.ID)

	if err := s.repo.UpdatePatientHISID(patient.ID, hisPatientID); err != nil {
		log.Printf("Failed to update HIS Patient ID: %v", err)
		return
	}

	s.hub.BroadcastPatientHISIDUpdate(patient.ID, hisPatientID)
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
	patient, err := s.repo.GetPatientByID(id)
	if err != nil {
		return errors.New("patient not found")
	}

	if err := s.repo.DeletePatient(id); err != nil {
		return err
	}

	s.hub.BroadcastPatientDeleted(id)

	if patient.HISPatientID != nil {
		go s.sendDeleteToHIS(id, *patient.HISPatientID)
	}

	return nil
}

func (s *PatientService) sendDeleteToHIS(patientID int, hisPatientID string) {
	messageID, hl7Message := hl7.GenerateADTA23(hisPatientID)

	log.Printf("Sending delete for patient %d to HIS via HL7 (MessageID: %s)", patientID, messageID)
	log.Printf("HL7 ADT^A23: %s", strings.ReplaceAll(string(hl7Message), "\r", "|"))

	ack, err := s.mllpClient.SendMessage(hl7Message)
	if err != nil {
		log.Printf("Failed to send HL7 delete message: %v", err)
		return
	}

	log.Printf("Received ACK: %s", strings.ReplaceAll(string(ack), "\r", "|"))
}
