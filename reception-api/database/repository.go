package database

import (
	"database/sql"
	"reception-api/models"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	db *sql.DB
	sq sq.StatementBuilderType
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Repository) CreateUser(username, passwordHash string) error {
	query := r.sq.Insert("users").
		Columns("username", "password_hash").
		Values(username, passwordHash)

	sqlRaw, args, _ := query.ToSql()
	_, err := r.db.Exec(sqlRaw, args...)
	return err
}

func (r *Repository) GetUserByUsername(username string) (*models.User, error) {
	query := r.sq.Select("id", "username", "password_hash", "created_at").
		From("users").
		Where(sq.Eq{"username": username})

	sqlRaw, args, _ := query.ToSql()
	row := r.db.QueryRow(sqlRaw, args...)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreatePatient(patient models.Patient) (int, error) {
	query := r.sq.Insert("patients").
		Columns("first_name", "last_name", "middle_name", "date_of_birth", "gender").
		Values(patient.FirstName, patient.LastName, patient.MiddleName, patient.DateOfBirth, patient.Gender).
		Suffix("RETURNING id")

	sqlRaw, args, _ := query.ToSql()
	var id int
	err := r.db.QueryRow(sqlRaw, args...).Scan(&id)
	return id, err
}

func (r *Repository) GetAllPatients() ([]models.Patient, error) {
	query := r.sq.Select("id", "his_patient_id", "first_name", "last_name", "middle_name", "date_of_birth", "gender", "created_at", "updated_at").
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
		err := rows.Scan(&p.ID, &p.HISPatientID, &p.FirstName, &p.LastName, &p.MiddleName, &p.DateOfBirth, &p.Gender, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		patients = append(patients, p)
	}

	return patients, nil
}

func (r *Repository) GetPatientByID(id int) (*models.Patient, error) {
	query := r.sq.Select("id", "his_patient_id", "first_name", "last_name", "middle_name", "date_of_birth", "gender", "created_at", "updated_at").
		From("patients").
		Where(sq.Eq{"id": id})

	sqlRaw, args, _ := query.ToSql()
	row := r.db.QueryRow(sqlRaw, args...)

	var p models.Patient
	err := row.Scan(&p.ID, &p.HISPatientID, &p.FirstName, &p.LastName, &p.MiddleName, &p.DateOfBirth, &p.Gender, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *Repository) DeletePatient(id int) error {
	query := r.sq.Delete("patients").Where(sq.Eq{"id": id})
	sqlRaw, args, _ := query.ToSql()
	_, err := r.db.Exec(sqlRaw, args...)
	return err
}

func (r *Repository) UpdatePatientHISID(patientID int, hisPatientID string) error {
	query := r.sq.Update("patients").
		Set("his_patient_id", hisPatientID).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": patientID})

	sqlRaw, args, _ := query.ToSql()
	_, err := r.db.Exec(sqlRaw, args...)
	return err
}

func (r *Repository) CreateHL7Message(msg models.HL7Message) error {
	query := r.sq.Insert("hl7_messages").
		Columns("message_id", "patient_id", "message_type", "status").
		Values(msg.MessageID, msg.PatientID, msg.MessageType, msg.Status)

	sqlRaw, args, _ := query.ToSql()
	_, err := r.db.Exec(sqlRaw, args...)
	return err
}

func (r *Repository) GetHL7MessageByMessageID(messageID string) (*models.HL7Message, error) {
	query := r.sq.Select("id", "message_id", "patient_id", "message_type", "status", "his_patient_id", "created_at", "ack_received_at").
		From("hl7_messages").
		Where(sq.Eq{"message_id": messageID})

	sqlRaw, args, _ := query.ToSql()
	row := r.db.QueryRow(sqlRaw, args...)

	var msg models.HL7Message
	err := row.Scan(&msg.ID, &msg.MessageID, &msg.PatientID, &msg.MessageType, &msg.Status, &msg.HISPatientID, &msg.CreatedAt, &msg.AckReceivedAt)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (r *Repository) UpdateHL7MessageStatus(messageID string, status string, hisPatientID string) error {
	now := time.Now()
	query := r.sq.Update("hl7_messages").
		Set("status", status).
		Set("his_patient_id", hisPatientID).
		Set("ack_received_at", now).
		Where(sq.Eq{"message_id": messageID})

	sqlRaw, args, _ := query.ToSql()
	_, err := r.db.Exec(sqlRaw, args...)
	return err
}
