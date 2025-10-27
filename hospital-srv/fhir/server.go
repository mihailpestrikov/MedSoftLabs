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
	"google.golang.org/protobuf/encoding/protojson"
)

type FHIRServer struct {
	practitionerService *services.PractitionerService
	encounterService    *services.EncounterService
}

func NewFHIRServer(practitionerService *services.PractitionerService, encounterService *services.EncounterService) *FHIRServer {
	return &FHIRServer{
		practitionerService: practitionerService,
		encounterService:    encounterService,
	}
}

func (s *FHIRServer) GetPractitioners(c *gin.Context) {
	practitioners, err := s.practitionerService.GetAllPractitioners()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var fhirPractitioners []map[string]interface{}
	for _, p := range practitioners {
		fhirResource := PractitionerToFHIR(p)
		jsonBytes, _ := protojson.Marshal(fhirResource)

		var resourceMap map[string]interface{}
		_ = json.Unmarshal(jsonBytes, &resourceMap)
		fhirPractitioners = append(fhirPractitioners, resourceMap)
	}

	c.JSON(http.StatusOK, gin.H{
		"resourceType": "Bundle",
		"type":         "searchset",
		"entry":        fhirPractitioners,
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

	c.JSON(http.StatusCreated, resourceMap)
}

func (s *FHIRServer) GetEncounters(c *gin.Context) {
	encounters, err := s.encounterService.GetAllEncounters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var fhirEncounters []map[string]interface{}
	for _, e := range encounters {
		fhirResource := EncounterToFHIR(e)
		jsonBytes, _ := protojson.Marshal(fhirResource)

		var resourceMap map[string]interface{}
		_ = json.Unmarshal(jsonBytes, &resourceMap)
		fhirEncounters = append(fhirEncounters, resourceMap)
	}

	c.JSON(http.StatusOK, gin.H{
		"resourceType": "Bundle",
		"type":         "searchset",
		"entry":        fhirEncounters,
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
