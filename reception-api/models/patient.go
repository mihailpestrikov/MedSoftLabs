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
