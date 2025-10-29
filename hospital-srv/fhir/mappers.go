package fhir

import (
	"fmt"
	"hospital-srv/models"
	"time"

	codespb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	encpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/encounter_go_proto"
	practpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/practitioner_go_proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PractitionerToFHIR(p models.Practitioner) *practpb.Practitioner {
	given := []*dtpb.String{{Value: p.FirstName}}
	if p.MiddleName != nil && *p.MiddleName != "" {
		given = append(given, &dtpb.String{Value: *p.MiddleName})
	}

	resource := &practpb.Practitioner{
		Id: &dtpb.Id{Value: p.ID},
		Name: []*dtpb.HumanName{
			{
				Family: &dtpb.String{Value: p.LastName},
				Given:  given,
			},
		},
		Qualification: []*practpb.Practitioner_Qualification{
			{
				Code: &dtpb.CodeableConcept{
					Text: &dtpb.String{Value: p.Specialization},
				},
			},
		},
	}

	return resource
}

func FHIRToPractitioner(fhirPrac *practpb.Practitioner) (models.Practitioner, error) {
	practitioner := models.Practitioner{}

	if fhirPrac.Id != nil {
		practitioner.ID = fhirPrac.Id.Value
	}

	if len(fhirPrac.Name) > 0 {
		name := fhirPrac.Name[0]
		if name.Family != nil {
			practitioner.LastName = name.Family.Value
		}
		if len(name.Given) > 0 && name.Given[0] != nil {
			practitioner.FirstName = name.Given[0].Value
		}
		if len(name.Given) > 1 && name.Given[1] != nil {
			middleName := name.Given[1].Value
			practitioner.MiddleName = &middleName
		}
	}

	if len(fhirPrac.Qualification) > 0 && fhirPrac.Qualification[0].Code != nil && fhirPrac.Qualification[0].Code.Text != nil {
		practitioner.Specialization = fhirPrac.Qualification[0].Code.Text.Value
	}

	return practitioner, nil
}

func EncounterToFHIR(e models.EncounterWithDetails) *encpb.Encounter {
	startTime := timestamppb.New(e.StartTime)

	patientDisplay := fmt.Sprintf("%s %s", e.Patient.LastName, e.Patient.FirstName)
	if e.Patient.MiddleName != nil && *e.Patient.MiddleName != "" {
		patientDisplay = fmt.Sprintf("%s %s %s", e.Patient.LastName, e.Patient.FirstName, *e.Patient.MiddleName)
	}

	practitionerDisplay := fmt.Sprintf("%s %s", e.Practitioner.LastName, e.Practitioner.FirstName)
	if e.Practitioner.MiddleName != nil && *e.Practitioner.MiddleName != "" {
		practitionerDisplay = fmt.Sprintf("%s %s %s", e.Practitioner.LastName, e.Practitioner.FirstName, *e.Practitioner.MiddleName)
	}
	practitionerDisplay = fmt.Sprintf("%s - %s", practitionerDisplay, e.Practitioner.Specialization)

	statusCode := codespb.EncounterStatusCode_ARRIVED
	switch e.Status {
	case "planned":
		statusCode = codespb.EncounterStatusCode_PLANNED
	case "arrived":
		statusCode = codespb.EncounterStatusCode_ARRIVED
	case "in-progress":
		statusCode = codespb.EncounterStatusCode_IN_PROGRESS
	case "completed":
		statusCode = codespb.EncounterStatusCode_FINISHED
	case "cancelled":
		statusCode = codespb.EncounterStatusCode_CANCELLED
	}

	resource := &encpb.Encounter{
		Id: &dtpb.Id{Value: e.ID},
		Status: &encpb.Encounter_StatusCode{
			Value: statusCode,
		},
		Subject: &dtpb.Reference{
			Reference: &dtpb.Reference_Uri{
				Uri: &dtpb.String{Value: fmt.Sprintf("Patient/%s", e.Patient.ID)},
			},
			Display: &dtpb.String{Value: patientDisplay},
		},
		Participant: []*encpb.Encounter_Participant{
			{
				Individual: &dtpb.Reference{
					Reference: &dtpb.Reference_Uri{
						Uri: &dtpb.String{Value: fmt.Sprintf("Practitioner/%s", e.Practitioner.ID)},
					},
					Display: &dtpb.String{Value: practitionerDisplay},
				},
			},
		},
		Period: &dtpb.Period{
			Start: &dtpb.DateTime{
				ValueUs:   startTime.AsTime().UnixMicro(),
				Precision: dtpb.DateTime_SECOND,
			},
		},
	}

	return resource
}

func FHIRToEncounter(fhirEnc *encpb.Encounter) (models.Encounter, error) {
	encounter := models.Encounter{
		Status: "arrived",
	}

	if fhirEnc.Id != nil {
		encounter.ID = fhirEnc.Id.Value
	}

	if fhirEnc.Status != nil {
		switch fhirEnc.Status.Value {
		case codespb.EncounterStatusCode_PLANNED:
			encounter.Status = "planned"
		case codespb.EncounterStatusCode_ARRIVED:
			encounter.Status = "arrived"
		case codespb.EncounterStatusCode_IN_PROGRESS:
			encounter.Status = "in-progress"
		case codespb.EncounterStatusCode_FINISHED:
			encounter.Status = "completed"
		case codespb.EncounterStatusCode_CANCELLED:
			encounter.Status = "cancelled"
		}
	}

	if fhirEnc.Subject != nil && fhirEnc.Subject.GetUri() != nil {
		var patientID string
		fmt.Sscanf(fhirEnc.Subject.GetUri().Value, "Patient/%s", &patientID)
		encounter.PatientID = patientID
	}

	if len(fhirEnc.Participant) > 0 && fhirEnc.Participant[0].Individual != nil {
		if uri := fhirEnc.Participant[0].Individual.GetUri(); uri != nil {
			var practitionerID string
			fmt.Sscanf(uri.Value, "Practitioner/%s", &practitionerID)
			encounter.PractitionerID = practitionerID
		}
	}

	if fhirEnc.Period != nil && fhirEnc.Period.Start != nil {
		encounter.StartTime = time.UnixMicro(fhirEnc.Period.Start.ValueUs)
	} else {
		encounter.StartTime = time.Now()
	}

	return encounter, nil
}
