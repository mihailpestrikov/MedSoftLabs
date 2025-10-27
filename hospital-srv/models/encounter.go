package models

import "time"

type Encounter struct {
	ID             string    `json:"id"`
	PatientID      string    `json:"patient_id"`
	PractitionerID string    `json:"practitioner_id"`
	Status         string    `json:"status"`
	StartTime      time.Time `json:"start_time"`
	CreatedAt      time.Time `json:"created_at"`
}

type EncounterWithDetails struct {
	Encounter
	Patient      Patient      `json:"patient"`
	Practitioner Practitioner `json:"practitioner"`
}
