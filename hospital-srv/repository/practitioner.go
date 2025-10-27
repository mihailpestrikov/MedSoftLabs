package repository

import (
	"hospital-srv/models"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) GetAllPractitioners() ([]models.Practitioner, error) {
	query := r.sq.Select("id", "first_name", "last_name", "middle_name", "specialization", "created_at").
		From("practitioners").
		OrderBy("last_name ASC")

	sqlRaw, args, _ := query.ToSql()
	rows, err := r.db.Query(sqlRaw, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var practitioners []models.Practitioner
	for rows.Next() {
		var p models.Practitioner
		err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.MiddleName, &p.Specialization, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		practitioners = append(practitioners, p)
	}

	return practitioners, nil
}

func (r *Repository) GetPractitionerByID(id string) (*models.Practitioner, error) {
	query := r.sq.Select("id", "first_name", "last_name", "middle_name", "specialization", "created_at").
		From("practitioners").
		Where(sq.Eq{"id": id})

	sqlRaw, args, _ := query.ToSql()
	row := r.db.QueryRow(sqlRaw, args...)

	var p models.Practitioner
	err := row.Scan(&p.ID, &p.FirstName, &p.LastName, &p.MiddleName, &p.Specialization, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
