package router

import (
	"hospital-srv/fhir"
	"hospital-srv/handlers"
	"hospital-srv/websocket"

	"github.com/gin-gonic/gin"
)

func Setup(patientHandler *handlers.PatientHandler, hub *websocket.Hub, fhirServer *fhir.FHIRServer) *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	})

	api := router.Group("/api")
	{
		patients := api.Group("/patients")
		{
			patients.GET("", patientHandler.GetAllPatients)
			patients.GET("/:id", patientHandler.GetPatient)
			patients.POST("", patientHandler.CreatePatient)
			patients.POST("/batch-delete", patientHandler.BatchDeletePatients)
			patients.DELETE("/:id", patientHandler.DeletePatient)
		}
	}

	fhirRoutes := router.Group("/fhir")
	{
		fhirRoutes.GET("/Practitioner", fhirServer.GetPractitioners)
		fhirRoutes.GET("/Practitioner/:id", fhirServer.GetPractitioner)
		fhirRoutes.POST("/Practitioner", fhirServer.CreatePractitioner)
		fhirRoutes.POST("/Encounter", fhirServer.CreateEncounter)
		fhirRoutes.GET("/Encounter", fhirServer.GetEncounters)
		fhirRoutes.GET("/Encounter/:id", fhirServer.GetEncounter)
	}

	return router
}
