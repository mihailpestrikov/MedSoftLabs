package fhir

import (
	"encoding/json"
	"hospital-srv/services"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	encpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/encounter_go_proto"
	practpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/practitioner_go_proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type FHIRServer struct {
	practitionerService *services.PractitionerService
	encounterService    *services.EncounterService
	notificationClient  *NotificationClient
}

func NewFHIRServer(practitionerService *services.PractitionerService, encounterService *services.EncounterService, notificationClient *NotificationClient) *FHIRServer {
	return &FHIRServer{
		practitionerService: practitionerService,
		encounterService:    encounterService,
		notificationClient:  notificationClient,
	}
}

func (s *FHIRServer) GetPractitioners(c *gin.Context) {
	practitioners, err := s.practitionerService.GetAllPractitioners()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var entries []map[string]interface{}
	for _, p := range practitioners {
		fhirResource := PractitionerToFHIR(p)
		jsonBytes, _ := protojson.Marshal(fhirResource)

		var resourceMap map[string]interface{}
		_ = json.Unmarshal(jsonBytes, &resourceMap)

		entry := map[string]interface{}{
			"resource": resourceMap,
		}
		entries = append(entries, entry)
	}

	c.JSON(http.StatusOK, gin.H{
		"resourceType": "Bundle",
		"type":         "searchset",
		"entry":        entries,
	})
}

func (s *FHIRServer) GetPractitioner(c *gin.Context) {
	id := c.Param("id")

	practitioner, err := s.practitionerService.GetPractitionerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Practitioner not found"})
		return
	}

	fhirResource := PractitionerToFHIR(*practitioner)
	jsonBytes, _ := protojson.Marshal(fhirResource)

	var resourceMap map[string]interface{}
	_ = json.Unmarshal(jsonBytes, &resourceMap)

	c.JSON(http.StatusOK, resourceMap)
}

func (s *FHIRServer) CreatePractitioner(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	log.Printf("Received FHIR Practitioner: %s", strings.ReplaceAll(string(body), "\n", " "))

	var fhirPractitioner practpb.Practitioner
	if err := protojson.Unmarshal(body, &fhirPractitioner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid FHIR Practitioner format"})
		return
	}

	practitioner, err := FHIRToPractitioner(&fhirPractitioner)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := s.practitionerService.CreatePractitioner(practitioner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createdPractitioner, err := s.practitionerService.GetPractitionerByID(id)
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"id": id})
		return
	}

	fhirResource := PractitionerToFHIR(*createdPractitioner)
	jsonBytes, _ := protojson.Marshal(fhirResource)

	log.Printf("Sending FHIR Practitioner response: %s", strings.ReplaceAll(string(jsonBytes), "\n", " "))

	var resourceMap map[string]interface{}
	_ = json.Unmarshal(jsonBytes, &resourceMap)

	c.JSON(http.StatusCreated, resourceMap)
}

func (s *FHIRServer) CreateEncounter(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	log.Printf("Received FHIR Encounter: %s", strings.ReplaceAll(string(body), "\n", " "))

	var fhirEncounter encpb.Encounter
	if err := protojson.Unmarshal(body, &fhirEncounter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid FHIR Encounter format"})
		return
	}

	encounter, err := FHIRToEncounter(&fhirEncounter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := s.encounterService.CreateEncounter(encounter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createdEncounter, err := s.encounterService.GetEncounterByID(id)
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"id": id})
		return
	}

	fhirResource := EncounterToFHIR(*createdEncounter)
	jsonBytes, _ := protojson.Marshal(fhirResource)

	log.Printf("Sending FHIR Encounter response: %s", strings.ReplaceAll(string(jsonBytes), "\n", " "))

	var resourceMap map[string]interface{}
	_ = json.Unmarshal(jsonBytes, &resourceMap)

	go func() {
		if err := s.notificationClient.NotifyEncounterCreated(resourceMap); err != nil {
		} else {
		}
	}()

	c.JSON(http.StatusCreated, resourceMap)
}

func (s *FHIRServer) GetEncounters(c *gin.Context) {
	encounters, err := s.encounterService.GetAllEncounters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var entries []map[string]interface{}
	for _, e := range encounters {
		fhirResource := EncounterToFHIR(e)
		jsonBytes, _ := protojson.Marshal(fhirResource)

		var resourceMap map[string]interface{}
		_ = json.Unmarshal(jsonBytes, &resourceMap)

		entry := map[string]interface{}{
			"resource": resourceMap,
		}
		entries = append(entries, entry)
	}

	c.JSON(http.StatusOK, gin.H{
		"resourceType": "Bundle",
		"type":         "searchset",
		"entry":        entries,
	})
}

func (s *FHIRServer) GetEncounter(c *gin.Context) {
	id := c.Param("id")

	encounter, err := s.encounterService.GetEncounterByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Encounter not found"})
		return
	}

	fhirResource := EncounterToFHIR(*encounter)
	jsonBytes, _ := protojson.Marshal(fhirResource)

	var resourceMap map[string]interface{}
	_ = json.Unmarshal(jsonBytes, &resourceMap)

	c.JSON(http.StatusOK, resourceMap)
}

func (s *FHIRServer) UpdateEncounterStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := s.encounterService.UpdateEncounterStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedEncounter, err := s.encounterService.GetEncounterByID(id)
	if err == nil {
		fhirResource := EncounterToFHIR(*updatedEncounter)
		jsonBytes, _ := protojson.Marshal(fhirResource)

		var resourceMap map[string]interface{}
		_ = json.Unmarshal(jsonBytes, &resourceMap)

		go func() {
			if err := s.notificationClient.NotifyEncounterStatusUpdated(resourceMap); err != nil {
				log.Printf("Failed to send encounter_status_updated notification: %v", err)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}
