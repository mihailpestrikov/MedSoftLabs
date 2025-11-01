package fhir

import (
	"doctor-api/models"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

// GetStringValue extracts string value from FHIR field which may be a string or map with "value" key.
func GetStringValue(field interface{}) string {
	if field == nil {
		return ""
	}
	if s, ok := field.(string); ok {
		return s
	}
	if m, ok := field.(map[string]interface{}); ok {
		if val, ok := m["value"].(string); ok {
			return val
		}
	}
	return ""
}

// GetInt64Value extracts int64 value from FHIR field, handling multiple formats.
func GetInt64Value(field interface{}) int64 {
	if field == nil {
		return 0
	}
	if i, ok := field.(int64); ok {
		return i
	}
	if f, ok := field.(float64); ok {
		return int64(f)
	}
	if s, ok := field.(string); ok {
		var val int64
		fmt.Sscanf(s, "%d", &val)
		return val
	}
	if m, ok := field.(map[string]interface{}); ok {
		if val, ok := m["value"].(int64); ok {
			return val
		}
		if val, ok := m["value"].(float64); ok {
			return int64(val)
		}
		if val, ok := m["value"].(string); ok {
			var i int64
			fmt.Sscanf(val, "%d", &i)
			return i
		}
	}
	return 0
}

func ExtractGenderFromDisplay(display string) string {
	re := regexp.MustCompile(`\[(male|female)\]$`)
	matches := re.FindStringSubmatch(display)
	if len(matches) > 1 {
		return strings.ToLower(matches[1])
	}
	return ""
}

func RemoveGenderFromDisplay(display string) string {
	re := regexp.MustCompile(`\s*\[(male|female)\]$`)
	return re.ReplaceAllString(display, "")
}

func ExtractIDFromReference(reference string) string {
	parts := strings.Split(reference, "/")
	if len(parts) == 2 {
		return parts[1]
	}
	return reference
}

func ParsePractitionerDisplay(display string) (string, string) {
	idx := strings.Index(display, " - ")
	if idx != -1 {
		name := display[:idx]
		specialization := display[idx+3:]
		return name, specialization
	}
	return display, ""
}

func NormalizeStatus(status string) string {
	status = strings.ToLower(status)
	status = strings.ReplaceAll(status, "_", "-")
	if status == "finished" {
		return "completed"
	}
	return status
}

// MapFHIRToEncounterDTO converts FHIR Encounter resource to EncounterDTO.
func MapFHIRToEncounterDTO(fhirData interface{}) (*models.EncounterDTO, error) {
	data, ok := fhirData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid FHIR data format")
	}

	dto := &models.EncounterDTO{}

	if id, ok := data["id"].(map[string]interface{}); ok {
		dto.ID = GetStringValue(id)
	}

	if status, ok := data["status"].(map[string]interface{}); ok {
		statusStr := GetStringValue(status)
		dto.Status = NormalizeStatus(statusStr)
	}

	if subject, ok := data["subject"].(map[string]interface{}); ok {
		if display, ok := subject["display"].(map[string]interface{}); ok {
			fullDisplay := GetStringValue(display)
			dto.PatientGender = ExtractGenderFromDisplay(fullDisplay)
			dto.PatientName = RemoveGenderFromDisplay(fullDisplay)
		}
		if reference, ok := subject["reference"].(map[string]interface{}); ok {
			ref := GetStringValue(reference)
			dto.PatientID = ExtractIDFromReference(ref)
		}
	}

	if participants, ok := data["participant"].([]interface{}); ok && len(participants) > 0 {
		if participant, ok := participants[0].(map[string]interface{}); ok {
			if individual, ok := participant["individual"].(map[string]interface{}); ok {
				if display, ok := individual["display"].(map[string]interface{}); ok {
					fullDisplay := GetStringValue(display)
					dto.PractitionerName, dto.PractitionerSpecialization = ParsePractitionerDisplay(fullDisplay)
				}
				if reference, ok := individual["reference"].(map[string]interface{}); ok {
					ref := GetStringValue(reference)
					dto.PractitionerID = ExtractIDFromReference(ref)
				}
			}
		}
	}

	if period, ok := data["period"].(map[string]interface{}); ok {
		if start, ok := period["start"].(map[string]interface{}); ok {
			valueUs := GetInt64Value(start["valueUs"])
			if valueUs > 0 {
				millis := valueUs / 1000
				t := time.UnixMilli(millis)
				dto.CreatedAt = t.Format(time.RFC3339)
			}
		}
	}

	if dto.CreatedAt == "" {
		dto.CreatedAt = time.Now().Format(time.RFC3339)
	}

	jsonBytes, _ := json.Marshal(dto)
	log.Printf("Mapped FHIR to DTO: %s", string(jsonBytes))

	return dto, nil
}

// MapFHIRToPractitionerDTO converts FHIR Practitioner resource to PractitionerDTO.
func MapFHIRToPractitionerDTO(fhirData interface{}) (*models.PractitionerDTO, error) {
	data, ok := fhirData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid FHIR data format")
	}

	dto := &models.PractitionerDTO{
		ID: GetStringValue(data["id"]),
	}

	if names, ok := data["name"].([]interface{}); ok && len(names) > 0 {
		if name, ok := names[0].(map[string]interface{}); ok {
			if family, ok := name["family"].(map[string]interface{}); ok {
				dto.LastName = GetStringValue(family)
			}
			if given, ok := name["given"].([]interface{}); ok {
				if len(given) > 0 {
					if g, ok := given[0].(map[string]interface{}); ok {
						dto.FirstName = GetStringValue(g)
					}
				}
				if len(given) > 1 {
					if g, ok := given[1].(map[string]interface{}); ok {
						dto.MiddleName = GetStringValue(g)
					}
				}
			}
		}
	}

	if qualifications, ok := data["qualification"].([]interface{}); ok && len(qualifications) > 0 {
		if q, ok := qualifications[0].(map[string]interface{}); ok {
			if code, ok := q["code"].(map[string]interface{}); ok {
				if text, ok := code["text"].(map[string]interface{}); ok {
					dto.Specialization = GetStringValue(text)
				}
			}
		}
	}

	return dto, nil
}
