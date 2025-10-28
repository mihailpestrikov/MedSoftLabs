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
