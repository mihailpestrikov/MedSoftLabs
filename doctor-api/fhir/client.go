package fhir

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"doctor-api/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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

func (c *FHIRClient) GetEncountersByPractitioner(practitionerID string) ([]models.EncounterDTO, error) {
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

		if !matchesPractitioner(resource, practitionerID) {
			continue
		}

		dto, err := MapFHIRToEncounterDTO(resource)
		if err != nil {
			log.Printf("Failed to map FHIR to DTO: %v", err)
			continue
		}

		encounters = append(encounters, *dto)
	}

	log.Printf("Returning %d encounters for practitioner %s", len(encounters), practitionerID)
	return encounters, nil
}

func (c *FHIRClient) UpdateEncounterStatus(encounterID string, status string) error {
	reqBody := map[string]string{
		"status": status,
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("Sending FHIR status update to HIS: encounter=%s, status=%s", encounterID, status)

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

	log.Printf("Successfully updated encounter status: encounter=%s, status=%s", encounterID, status)
	return nil
}

func matchesPractitioner(resource map[string]interface{}, practitionerID string) bool {
	participants, ok := resource["participant"].([]interface{})
	if !ok || len(participants) == 0 {
		return false
	}

	for _, p := range participants {
		participant, ok := p.(map[string]interface{})
		if !ok {
			continue
		}

		individual, ok := participant["individual"].(map[string]interface{})
		if !ok {
			continue
		}

		reference, ok := individual["reference"].(map[string]interface{})
		if !ok {
			continue
		}

		refValue, ok := reference["value"].(string)
		if !ok {
			continue
		}

		if strings.HasSuffix(refValue, practitionerID) {
			return true
		}
	}

	return false
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

type Practitioner struct {
	ID             string  `json:"id"`
	FirstName      string  `json:"firstName"`
	LastName       string  `json:"lastName"`
	MiddleName     *string `json:"middleName,omitempty"`
	Specialization string  `json:"specialization"`
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

	log.Printf("Received FHIR Practitioner response from HIS")

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HIS returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var bundle map[string]interface{}
	if err := json.Unmarshal(respBody, &bundle); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	entries, ok := bundle["entry"].([]interface{})
	if !ok {
		return []Practitioner{}, nil
	}

	var practitioners []Practitioner
	for _, entry := range entries {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		resource, ok := entryMap["resource"].(map[string]interface{})
		if !ok {
			continue
		}

		practitioner := Practitioner{
			ID: fmt.Sprintf("%v", getValue(resource, "id")),
		}

		if names, ok := resource["name"].([]interface{}); ok && len(names) > 0 {
			if name, ok := names[0].(map[string]interface{}); ok {
				if family, ok := name["family"].(map[string]interface{}); ok {
					practitioner.LastName = fmt.Sprintf("%v", family["value"])
				}
				if given, ok := name["given"].([]interface{}); ok {
					if len(given) > 0 {
						if g, ok := given[0].(map[string]interface{}); ok {
							practitioner.FirstName = fmt.Sprintf("%v", g["value"])
						}
					}
					if len(given) > 1 {
						if g, ok := given[1].(map[string]interface{}); ok {
							middleName := fmt.Sprintf("%v", g["value"])
							practitioner.MiddleName = &middleName
						}
					}
				}
			}
		}

		if qualifications, ok := resource["qualification"].([]interface{}); ok && len(qualifications) > 0 {
			if q, ok := qualifications[0].(map[string]interface{}); ok {
				if code, ok := q["code"].(map[string]interface{}); ok {
					if text, ok := code["text"].(map[string]interface{}); ok {
						practitioner.Specialization = fmt.Sprintf("%v", text["value"])
					}
				}
			}
		}

		practitioners = append(practitioners, practitioner)
	}

	return practitioners, nil
}
