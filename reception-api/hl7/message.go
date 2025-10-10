package hl7

import (
	"fmt"
	"reception-api/models"
	"strings"
	"time"

	"github.com/google/simhospital/pkg/hl7"
	"github.com/google/uuid"
)

func init() {
	hl7.TimezoneAndLocation("UTC")
}

func GenerateADTA04(patient *models.Patient) (string, []byte) {
	timestamp := time.Now().Format("20060102150405")
	messageID := uuid.New().String()

	msh := fmt.Sprintf("MSH|^~\\&|RECEPTION|CLINIC|HIS|HOSPITAL|%s||ADT^A04|%s|P|2.5", timestamp, messageID)

	dob := strings.ReplaceAll(patient.DateOfBirth, "-", "")

	pid := fmt.Sprintf("PID|||%d||%s^%s^%s||%s|%s",
		patient.ID,
		patient.LastName,
		patient.FirstName,
		valueOrEmpty(patient.MiddleName),
		dob,
		strings.ToUpper(patient.Gender))

	message := fmt.Sprintf("%s\r%s", msh, pid)
	return messageID, []byte(message)
}

func GenerateADTA23(hisPatientID string) (string, []byte) {
	timestamp := time.Now().Format("20060102150405")
	messageID := uuid.New().String()

	msh := fmt.Sprintf("MSH|^~\\&|RECEPTION|CLINIC|HIS|HOSPITAL|%s||ADT^A23|%s|P|2.5", timestamp, messageID)
	pid := fmt.Sprintf("PID|||%s", hisPatientID)

	message := fmt.Sprintf("%s\r%s", msh, pid)
	return messageID, []byte(message)
}

func ParseACK(data []byte) (string, string, error) {
	msg, err := hl7.ParseMessage(data)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse ACK: %w", err)
	}

	var originalMessageID string
	var patientUUID string

	msas, err := msg.AllMSA()
	if err == nil && len(msas) > 0 {
		msa := msas[0]
		if msa.MessageControlID != nil {
			originalMessageID = string(*msa.MessageControlID)
		}
	}

	pids, err := msg.AllPID()
	if err == nil && len(pids) > 0 {
		pid := pids[0]
		if len(pid.PatientIdentifierList) > 0 && pid.PatientIdentifierList[0].IDNumber != nil {
			patientUUID = string(*pid.PatientIdentifierList[0].IDNumber)
		}
	}

	if originalMessageID == "" {
		return "", "", fmt.Errorf("no message ID in ACK")
	}

	return originalMessageID, patientUUID, nil
}

func valueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
