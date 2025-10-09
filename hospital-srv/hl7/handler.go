package hl7

import (
	"hospital-srv/models"
	"hospital-srv/services"
	"log"
	"strings"
)

type HL7Handler struct {
	patientService *services.PatientService
}

func NewHL7Handler(patientService *services.PatientService) *HL7Handler {
	return &HL7Handler{
		patientService: patientService,
	}
}

func (h *HL7Handler) HandleMessage(data []byte) []byte {
	log.Printf("Received HL7 message: %s", strings.ReplaceAll(string(data), "\r", "|"))

	msg, err := ParseHL7(data)
	if err != nil {
		log.Printf("Error parsing HL7 message: %v", err)
		return []byte("MSH|^~\\&|HIS|HOSPITAL|||" + "|ACK^A04||P|2.5\rMSA|AE|UNKNOWN")
	}

	switch msg.MessageType {
	case "ADT^A04":
		return h.handlePatientAdmit(msg)
	case "ADT^A23":
		return h.handlePatientDelete(msg)
	default:
		log.Printf("Unknown message type: %s", msg.MessageType)
		return []byte("MSH|^~\\&|HIS|HOSPITAL|||" + "|ACK|P|2.5\rMSA|AR|" + msg.MessageID)
	}
}

func (h *HL7Handler) handlePatientAdmit(msg *HL7Message) []byte {
	patient := models.Patient{
		FirstName:   msg.FirstName,
		LastName:    msg.LastName,
		MiddleName:  &msg.MiddleName,
		DateOfBirth: msg.DateOfBirth,
		Gender:      strings.ToLower(msg.Gender),
	}

	uuid, err := h.patientService.CreatePatient(patient)
	if err != nil {
		log.Printf("Error creating patient: %v", err)
		return []byte("MSH|^~\\&|HIS|HOSPITAL|||" + "|ACK^A04||P|2.5\rMSA|AE|" + msg.MessageID)
	}

	log.Printf("Created patient with UUID: %s", uuid)

	ack := GenerateACK(msg.MessageID, uuid)
	log.Printf("Sending ACK: %s", strings.ReplaceAll(string(ack), "\r", "|"))
	return ack
}

func (h *HL7Handler) handlePatientDelete(msg *HL7Message) []byte {
	err := h.patientService.DeletePatient(msg.PatientID)
	if err != nil {
		log.Printf("Error deleting patient: %v", err)
		return []byte("MSH|^~\\&|HIS|HOSPITAL|||" + "|ACK^A23||P|2.5\rMSA|AE|" + msg.MessageID)
	}

	log.Printf("Deleted patient: %s", msg.PatientID)

	ack := GenerateACK(msg.MessageID, msg.PatientID)
	log.Printf("Sending ACK: %s", strings.ReplaceAll(string(ack), "\r", "|"))
	return ack
}
