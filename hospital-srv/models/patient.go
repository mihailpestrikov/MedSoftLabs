package models

import "time"

type Patient struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	MiddleName  *string   `json:"middle_name"`
	DateOfBirth string    `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
