package models

// EncounterDTO represents encounter data for client applications.
type EncounterDTO struct {
	ID                         string `json:"id"`
	PatientID                  string `json:"patientId"`
	PatientName                string `json:"patientName"`
	PatientGender              string `json:"patientGender"`
	PractitionerID             string `json:"practitionerId"`
	PractitionerName           string `json:"practitionerName"`
	PractitionerSpecialization string `json:"practitionerSpecialization"`
	Status                     string `json:"status"`
	CreatedAt                  string `json:"createdAt"`
}

// PractitionerDTO represents practitioner data for client applications.
type PractitionerDTO struct {
	ID             string `json:"id"`
	FirstName      string `json:"firstName"`
	MiddleName     string `json:"middleName,omitempty"`
	LastName       string `json:"lastName"`
	Specialization string `json:"specialization"`
}
