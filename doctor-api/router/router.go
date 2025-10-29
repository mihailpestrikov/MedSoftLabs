package router

import (
	"doctor-api/handlers"
	"doctor-api/websocket"

	"github.com/gin-gonic/gin"
)

func Setup(encounterHandler *handlers.EncounterHandler, practitionerHandler *handlers.PractitionerHandler, fhirNotificationHandler *handlers.FHIRNotificationHandler, hub *websocket.Hub) *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	})

	fhir := router.Group("/fhir")
	{
		fhir.POST("/notifications/encounter", fhirNotificationHandler.HandleEncounterNotification)
	}

	api := router.Group("/api")
	{
		api.GET("/practitioners", practitionerHandler.GetAllPractitioners)
		api.GET("/encounters/:practitioner_id", encounterHandler.GetEncountersByPractitioner)
		api.PATCH("/encounters/:id", encounterHandler.UpdateEncounterStatus)
	}

	return router
}
