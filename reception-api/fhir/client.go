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
	"reception-api/models"
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
			Value: codespb.EncounterStatusCode_PLANNED,
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

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("HIS returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	encounterID, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("no id in response")
	}

	log.Printf("Created encounter with ID: %s", encounterID)

	return encounterID, nil
}

func (c *FHIRClient) GetPractitioners() ([]models.PractitionerDTO, error) {
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

	var bundle map[string]interface{}
	if err := json.Unmarshal(respBody, &bundle); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	entries, ok := bundle["entry"].([]interface{})
	if !ok {
		return []models.PractitionerDTO{}, nil
	}

	var practitioners []models.PractitionerDTO
	for _, entry := range entries {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		resource, ok := entryMap["resource"].(map[string]interface{})
		if !ok {
			continue
		}

		dto, err := MapFHIRToPractitionerDTO(resource)
		if err != nil {
			log.Printf("Failed to map FHIR Practitioner to DTO: %v", err)
			continue
		}

		practitioners = append(practitioners, *dto)
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

func (c *FHIRClient) GetEncounters() ([]models.EncounterDTO, error) {
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
		return []models.EncounterDTO{}, nil
	}

	var encounters []models.EncounterDTO
	for _, entry := range entries {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		resource, ok := entryMap["resource"].(map[string]interface{})
		if !ok {
			continue
		}

		dto, err := MapFHIRToEncounterDTO(resource)
		if err != nil {
			log.Printf("Failed to map FHIR to DTO: %v", err)
			continue
		}

		encounters = append(encounters, *dto)
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
