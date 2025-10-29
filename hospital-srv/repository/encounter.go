package repository

import (
	"hospital-srv/models"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) CreateEncounter(encounter models.Encounter) (string, error) {
	query := r.sq.Insert("encounters").
		Columns("patient_id", "practitioner_id", "status", "start_time").
		Values(encounter.PatientID, encounter.PractitionerID, encounter.Status, encounter.StartTime).
		Suffix("RETURNING id")

	sqlRaw, args, _ := query.ToSql()
	var id string
	err := r.db.QueryRow(sqlRaw, args...).Scan(&id)
	return id, err
}

func (r *Repository) GetAllEncounters() ([]models.EncounterWithDetails, error) {
	query := r.sq.Select(
		"e.id", "e.patient_id", "e.practitioner_id", "e.status", "e.start_time", "e.created_at",
		"pat.id", "pat.first_name", "pat.last_name", "pat.middle_name", "pat.date_of_birth", "pat.gender", "pat.created_at", "pat.updated_at",
		"pr.id", "pr.first_name", "pr.last_name", "pr.middle_name", "pr.specialization", "pr.created_at",
	).
		From("encounters e").
		Join("patients pat ON e.patient_id = pat.id").
		Join("practitioners pr ON e.practitioner_id = pr.id").
		OrderBy("e.start_time DESC")

	sqlRaw, args, _ := query.ToSql()
	rows, err := r.db.Query(sqlRaw, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var encounters []models.EncounterWithDetails
	for rows.Next() {
		var e models.EncounterWithDetails
		err := rows.Scan(
			&e.ID, &e.PatientID, &e.PractitionerID, &e.Status, &e.StartTime, &e.CreatedAt,
			&e.Patient.ID, &e.Patient.FirstName, &e.Patient.LastName, &e.Patient.MiddleName, &e.Patient.DateOfBirth, &e.Patient.Gender, &e.Patient.CreatedAt, &e.Patient.UpdatedAt,
			&e.Practitioner.ID, &e.Practitioner.FirstName, &e.Practitioner.LastName, &e.Practitioner.MiddleName, &e.Practitioner.Specialization, &e.Practitioner.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		encounters = append(encounters, e)
	}

	return encounters, nil
}

func (r *Repository) GetEncounterByID(id string) (*models.EncounterWithDetails, error) {
	query := r.sq.Select(
		"e.id", "e.patient_id", "e.practitioner_id", "e.status", "e.start_time", "e.created_at",
		"pat.id", "pat.first_name", "pat.last_name", "pat.middle_name", "pat.date_of_birth", "pat.gender", "pat.created_at", "pat.updated_at",
		"pr.id", "pr.first_name", "pr.last_name", "pr.middle_name", "pr.specialization", "pr.created_at",
	).
		From("encounters e").
		Join("patients pat ON e.patient_id = pat.id").
		Join("practitioners pr ON e.practitioner_id = pr.id").
		Where(sq.Eq{"e.id": id})

	sqlRaw, args, _ := query.ToSql()
	row := r.db.QueryRow(sqlRaw, args...)

	var e models.EncounterWithDetails
	err := row.Scan(
		&e.ID, &e.PatientID, &e.PractitionerID, &e.Status, &e.StartTime, &e.CreatedAt,
		&e.Patient.ID, &e.Patient.FirstName, &e.Patient.LastName, &e.Patient.MiddleName, &e.Patient.DateOfBirth, &e.Patient.Gender, &e.Patient.CreatedAt, &e.Patient.UpdatedAt,
		&e.Practitioner.ID, &e.Practitioner.FirstName, &e.Practitioner.LastName, &e.Practitioner.MiddleName, &e.Practitioner.Specialization, &e.Practitioner.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *Repository) UpdateEncounterStatus(id string, status string) error {
	query := r.sq.Update("encounters").
		Set("status", status).
		Where(sq.Eq{"id": id})

	sqlRaw, args, _ := query.ToSql()
	_, err := r.db.Exec(sqlRaw, args...)
	return err
}
