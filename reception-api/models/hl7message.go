package models

import "time"

type HL7Message struct {
	ID            int        `json:"id"`
	MessageID     string     `json:"message_id"`
	PatientID     int        `json:"patient_id"`
	MessageType   string     `json:"message_type"`
	Status        string     `json:"status"`
	HISPatientID  *string    `json:"his_patient_id"`
	CreatedAt     time.Time  `json:"created_at"`
	AckReceivedAt *time.Time `json:"ack_received_at"`
}
