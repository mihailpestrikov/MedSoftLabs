package hl7

import (
	"fmt"
	"time"

	"github.com/google/simhospital/pkg/hl7"
)

func init() {
	hl7.TimezoneAndLocation("UTC")
}

type HL7Message struct {
	MessageType string
	MessageID   string
	PatientID   string
	FirstName   string
	LastName    string
	MiddleName  string
	DateOfBirth string
	Gender      string
}

func ParseHL7(data []byte) (*HL7Message, error) {
	msg, err := hl7.ParseMessage(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HL7 message: %w", err)
	}

	result := &HL7Message{}

	msh, err := msg.MSH()
	if err == nil && msh != nil {
		if msh.MessageType != nil {
			result.MessageType = fmt.Sprintf("%s^%s", msh.MessageType.MessageCode, msh.MessageType.TriggerEvent)
		}
		if msh.MessageControlID != nil {
			result.MessageID = string(*msh.MessageControlID)
		}
	}

	pids, err := msg.AllPID()
	if err == nil && len(pids) > 0 {
		pid := pids[0]

		if len(pid.PatientIdentifierList) > 0 {
			result.PatientID = string(*pid.PatientIdentifierList[0].IDNumber)
		}

		if len(pid.PatientName) > 0 {
			name := pid.PatientName[0]
			if name.FamilyName != nil && name.FamilyName.Surname != nil {
				result.LastName = string(*name.FamilyName.Surname)
			}
			if name.GivenName != nil {
				result.FirstName = string(*name.GivenName)
			}
			if name.SecondAndFurtherGivenNamesOrInitialsThereof != nil {
				result.MiddleName = string(*name.SecondAndFurtherGivenNamesOrInitialsThereof)
			}
		}

		if pid.DateTimeOfBirth != nil && !pid.DateTimeOfBirth.IsHL7Null {
			result.DateOfBirth = pid.DateTimeOfBirth.Time.Format("20060102")
		}

		if pid.AdministrativeSex != nil {
			result.Gender = string(*pid.AdministrativeSex)
		}
	}

	return result, nil
}

func GenerateACK(originalMessageID string, patientUUID string) []byte {
	timestamp := time.Now().Format("20060102150405")

	msh := fmt.Sprintf("MSH|^~\\&|HIS|HOSPITAL|RECEPTION|CLINIC|%s||ACK^A04|%s|P|2.5", timestamp, timestamp)
	msa := fmt.Sprintf("MSA|AA|%s", originalMessageID)
	pid := fmt.Sprintf("PID|||%s", patientUUID)

	ack := fmt.Sprintf("%s\r%s\r%s", msh, msa, pid)
	return []byte(ack)
}
