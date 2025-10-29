package fhir

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	codespb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	dtpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	encpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/encounter_go_proto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FHIRClient struct {
	baseURL    string
	httpClient *http.Client
}

type Practitioner struct {
	ID             string
	FirstName      string
	LastName       string
	MiddleName     string
	Specialization string
}

type FHIRBundle struct {
	Entry []FHIRPractitionerEntry `json:"entry"`
}

type FHIRPractitionerEntry struct {
	Resource FHIRPractitionerResource `json:"resource"`
}

type FHIRPractitionerResource struct {
	ID            FHIRValue           `json:"id"`
	Name          []FHIRName          `json:"name"`
	Qualification []FHIRQualification `json:"qualification"`
}

type FHIRValue struct {
	Value string `json:"value"`
}

type FHIRName struct {
	Family FHIRValue   `json:"family"`
	Given  []FHIRValue `json:"given"`
}

type FHIRQualification struct {
	Code FHIRCode `json:"code"`
}

type FHIRCode struct {
	Text FHIRValue `json:"text"`
}

func NewFHIRClient(baseURL string, certPath string) (*FHIRClient, error) {
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		return nil, fmt.Errorf("failed to append certificate")
	}

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return &FHIRClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}, nil
}

func (c *FHIRClient) CreateEncounter(patientID string, practitionerID string, startTime time.Time) (string, error) {
	startTimestamp := timestamppb.New(startTime)

	encounter := &encpb.Encounter{
		Status: &encpb.Encounter_StatusCode{
			Value: codespb.EncounterStatusCode_ARRIVED,
		},
		Subject: &dtpb.Reference{
			Reference: &dtpb.Reference_Uri{
				Uri: &dtpb.String{Value: fmt.Sprintf("Patient/%s", patientID)},
			},
		},
		Participant: []*encpb.Encounter_Participant{
			{
				Individual: &dtpb.Reference{
					Reference: &dtpb.Reference_Uri{
						Uri: &dtpb.String{Value: fmt.Sprintf("Practitioner/%s", practitionerID)},
					},
				},
			},
		},
		Period: &dtpb.Period{
			Start: &dtpb.DateTime{
				ValueUs:   startTimestamp.AsTime().UnixMicro(),
				Precision: dtpb.DateTime_SECOND,
			},
		},
	}

	jsonBytes, err := protojson.Marshal(encounter)
	if err != nil {
		return "", fmt.Errorf("failed to marshal encounter: %w", err)
	}

	log.Printf("Sending FHIR Encounter to HIS: %s", strings.ReplaceAll(string(jsonBytes), "\n", " "))

	url := fmt.Sprintf("%s/fhir/Encounter", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("Received FHIR response from HIS: %s", strings.ReplaceAll(string(respBody), "\n", " "))

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("HIS returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	idField, ok := result["id"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("no id in response")
	}

	encounterID, ok := idField["value"].(string)
	if !ok {
		return "", fmt.Errorf("invalid id format in response")
	}

	return encounterID, nil
}

func (c *FHIRClient) GetPractitioners() ([]Practitioner, error) {
	url := fmt.Sprintf("%s/fhir/Practitioner", c.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HIS returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var bundle FHIRBundle
	if err := json.Unmarshal(respBody, &bundle); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var practitioners []Practitioner
	for _, entry := range bundle.Entry {
		resource := entry.Resource
		firstName := ""
		lastName := ""
		middleName := ""
		specialization := ""

		if len(resource.Name) > 0 {
			lastName = resource.Name[0].Family.Value
			if len(resource.Name[0].Given) > 0 {
				firstName = resource.Name[0].Given[0].Value
			}
			if len(resource.Name[0].Given) > 1 {
				middleName = resource.Name[0].Given[1].Value
			}
		}

		if len(resource.Qualification) > 0 {
			specialization = resource.Qualification[0].Code.Text.Value
		}

		practitioners = append(practitioners, Practitioner{
			ID:             resource.ID.Value,
			FirstName:      firstName,
			LastName:       lastName,
			MiddleName:     middleName,
			Specialization: specialization,
		})
	}

	return practitioners, nil
}

func (c *FHIRClient) UpdateEncounterStatus(encounterID string, status string) error {
	reqBody := map[string]string{
		"status": status,
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/fhir/Encounter/%s", c.baseURL, encounterID)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HIS returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

type Encounter struct {
	ID             string
	PatientID      int
	PractitionerID string
	Status         string
	StartTime      string
	Patient        interface{}
	Practitioner   interface{}
}

func (c *FHIRClient) GetEncounters() ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/fhir/Encounter", c.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HIS returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var bundle map[string]interface{}
	if err := json.Unmarshal(respBody, &bundle); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	entries, ok := bundle["entry"].([]interface{})
	if !ok {
		return []map[string]interface{}{}, nil
	}

	var encounters []map[string]interface{}
	for _, entry := range entries {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		resource, ok := entryMap["resource"].(map[string]interface{})
		if !ok {
			continue
		}

		statusValue := strings.ToLower(fmt.Sprintf("%v", getValue(resource, "status")))
		statusValue = strings.ReplaceAll(statusValue, "_", "-")
		if statusValue == "finished" {
			statusValue = "completed"
		}

		encounter := map[string]interface{}{
			"id":     getValue(resource, "id"),
			"status": statusValue,
		}

		if period, ok := resource["period"].(map[string]interface{}); ok {
			if start, ok := period["start"].(map[string]interface{}); ok {
				if valueUs, ok := start["valueUs"].(string); ok {
					var microseconds int64
					fmt.Sscanf(valueUs, "%d", &microseconds)
					milliseconds := microseconds / 1000
					encounter["start_time"] = time.UnixMilli(milliseconds).Format(time.RFC3339)
				}
			}
		}

		if subject, ok := resource["subject"].(map[string]interface{}); ok {
			patientDisplay := getDisplayValue(subject)
			parts := strings.Fields(patientDisplay)
			patient := map[string]interface{}{}
			if len(parts) >= 2 {
				patient["last_name"] = parts[0]
				patient["first_name"] = parts[1]
				if len(parts) >= 3 {
					patient["middle_name"] = parts[2]
				}
			}
			encounter["patient"] = patient
		}

		if participants, ok := resource["participant"].([]interface{}); ok && len(participants) > 0 {
			if participant, ok := participants[0].(map[string]interface{}); ok {
				if individual, ok := participant["individual"].(map[string]interface{}); ok {
					practDisplay := getDisplayValue(individual)
					practitioner := map[string]interface{}{}
					if idx := strings.Index(practDisplay, " - "); idx != -1 {
						namePart := practDisplay[:idx]
						specialization := practDisplay[idx+3:]
						parts := strings.Fields(namePart)
						if len(parts) >= 2 {
							practitioner["LastName"] = parts[0]
							practitioner["FirstName"] = parts[1]
							if len(parts) >= 3 {
								practitioner["MiddleName"] = parts[2]
							}
							practitioner["Specialization"] = specialization
						}
					}
					encounter["practitioner"] = practitioner
				}
			}
		}

		encounters = append(encounters, encounter)
	}

	return encounters, nil
}

func getValue(m map[string]interface{}, key string) interface{} {
	if v, ok := m[key]; ok {
		if vm, ok := v.(map[string]interface{}); ok {
			if val, ok := vm["value"]; ok {
				return val
			}
		}
		return v
	}
	return nil
}

func getDisplayValue(ref map[string]interface{}) string {
	if display, ok := ref["display"].(map[string]interface{}); ok {
		if val, ok := display["value"].(string); ok {
			return val
		}
	}
	return ""
}
