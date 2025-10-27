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

func EncounterToFHIR(e models.EncounterWithDetails) *encpb.Encounter {
	startTime := timestamppb.New(e.StartTime)

	resource := &encpb.Encounter{
		Id: &dtpb.Id{Value: e.ID},
		Status: &encpb.Encounter_StatusCode{
			Value: codespb.EncounterStatusCode_ARRIVED,
		},
		Subject: &dtpb.Reference{
			Reference: &dtpb.Reference_Uri{
				Uri: &dtpb.String{Value: fmt.Sprintf("Patient/%s", e.PatientID)},
			},
		},
		Participant: []*encpb.Encounter_Participant{
			{
				Individual: &dtpb.Reference{
					Reference: &dtpb.Reference_Uri{
						Uri: &dtpb.String{Value: fmt.Sprintf("Practitioner/%s", e.PractitionerID)},
					},
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
