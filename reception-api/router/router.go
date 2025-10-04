package router

import (
	"reception-api/handlers"
	"reception-api/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(authHandler *handlers.AuthHandler, patientHandler *handlers.PatientHandler, jwtService *middleware.JWTService) *gin.Engine {
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
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}

		patients := api.Group("/patients")
		patients.Use(middleware.AuthMiddleware(jwtService))
		{
			patients.GET("", patientHandler.GetAllPatients)
			patients.POST("", patientHandler.CreatePatient)
			patients.GET("/:id", patientHandler.GetPatient)
			patients.DELETE("/:id", patientHandler.DeletePatient)
		}
	}

	return router
}
