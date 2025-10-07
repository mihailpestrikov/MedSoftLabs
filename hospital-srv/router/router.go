package router

import (
	"hospital-srv/handlers"

	"github.com/gin-gonic/gin"
)

func Setup(patientHandler *handlers.PatientHandler) *gin.Engine {
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

	api := router.Group("/api")
	{
		patients := api.Group("/patients")
		{
			patients.GET("", patientHandler.GetAllPatients)
			patients.GET("/:id", patientHandler.GetPatient)
			patients.POST("", patientHandler.CreatePatient)
			patients.DELETE("/:id", patientHandler.DeletePatient)
		}
	}

	return router
}
