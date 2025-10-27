package repository

import (
	"hospital-srv/models"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) CreatePatient(patient models.Patient) (string, error) {
	query := r.sq.Insert("patients").
		Columns("first_name", "last_name", "middle_name", "date_of_birth", "gender").
		Values(patient.FirstName, patient.LastName, patient.MiddleName, patient.DateOfBirth, patient.Gender).
		Suffix("RETURNING id")

	sqlRaw, args, _ := query.ToSql()
	var id string
	err := r.db.QueryRow(sqlRaw, args...).Scan(&id)
	return id, err
}

func (r *Repository) GetAllPatients() ([]models.Patient, error) {
	query := r.sq.Select("id", "first_name", "last_name", "middle_name", "date_of_birth", "gender", "created_at", "updated_at").
		From("patients").
		OrderBy("created_at DESC")

	sqlRaw, args, _ := query.ToSql()
	rows, err := r.db.Query(sqlRaw, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []models.Patient
	for rows.Next() {
		var p models.Patient
		err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.MiddleName, &p.DateOfBirth, &p.Gender, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		patients = append(patients, p)
	}

	return patients, nil
}

func (r *Repository) GetPatientByID(id string) (*models.Patient, error) {
	query := r.sq.Select("id", "first_name", "last_name", "middle_name", "date_of_birth", "gender", "created_at", "updated_at").
		From("patients").
		Where(sq.Eq{"id": id})

	sqlRaw, args, _ := query.ToSql()
	row := r.db.QueryRow(sqlRaw, args...)

	var p models.Patient
	err := row.Scan(&p.ID, &p.FirstName, &p.LastName, &p.MiddleName, &p.DateOfBirth, &p.Gender, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *Repository) DeletePatient(id string) error {
	query := r.sq.Delete("patients").Where(sq.Eq{"id": id})
	sqlRaw, args, _ := query.ToSql()
	_, err := r.db.Exec(sqlRaw, args...)
	return err
}

func (r *Repository) DeletePatients(ids []string) error {
	query := r.sq.Delete("patients").Where(sq.Eq{"id": ids})
	sqlRaw, args, _ := query.ToSql()
	_, err := r.db.Exec(sqlRaw, args...)
	return err
}
