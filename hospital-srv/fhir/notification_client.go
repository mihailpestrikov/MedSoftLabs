package fhir

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type NotificationClient struct {
	doctorAPIURL    string
	receptionAPIURL string
	httpClient      *http.Client
}

func NewNotificationClient(doctorAPIURL string, receptionAPIURL string) *NotificationClient {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return &NotificationClient{
		doctorAPIURL:    doctorAPIURL,
		receptionAPIURL: receptionAPIURL,
		httpClient:      httpClient,
	}
}

func (c *NotificationClient) NotifyEncounterCreated(encounterData interface{}) error {
	notification := map[string]interface{}{
		"type": "encounter_created",
		"data": encounterData,
	}

	c.sendNotificationToAll(notification)
	return nil
}

func (c *NotificationClient) NotifyEncounterStatusUpdated(encounterData interface{}) error {
	notification := map[string]interface{}{
		"type": "encounter_status_updated",
		"data": encounterData,
	}

	c.sendNotificationToAll(notification)
	return nil
}

func (c *NotificationClient) sendNotificationToAll(notification map[string]interface{}) {
	go c.sendNotification(c.doctorAPIURL, notification)
	go c.sendNotification(c.receptionAPIURL, notification)
}

func (c *NotificationClient) sendNotification(targetAPIURL string, notification map[string]interface{}) error {
	jsonBytes, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	url := fmt.Sprintf("%s/fhir/notifications/encounter", targetAPIURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	log.Printf("Sending FHIR notification: type=%s, url=%s", notification["type"], url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to send notification to %s: %v", targetAPIURL, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		log.Printf("API %s returned error: status=%d, body=%s", targetAPIURL, resp.StatusCode, string(respBody))
		return fmt.Errorf("notification failed with status %d", resp.StatusCode)
	}

	log.Printf("Successfully sent FHIR notification to %s: type=%s", targetAPIURL, notification["type"])
	return nil
}
