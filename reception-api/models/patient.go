package models

import "time"

type Patient struct {
	ID           int       `json:"id"`
	HISPatientID *string   `json:"his_patient_id"`
	FirstName    string    `json:"first_name" binding:"required"`
	LastName     string    `json:"last_name" binding:"required"`
	DateOfBirth  string    `json:"date_of_birth" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

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
