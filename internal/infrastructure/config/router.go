package config

import (
	"FleetManagementSystem/internal/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(locationHandler *controller.LocationController) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	apiV1 := r.Group("/api/v1")

	vehicles := apiV1.Group("/vehicles")

	{
		vehicles.GET("/:vehicle_id/location", locationHandler.GetLatestLocation)
		vehicles.GET("/:vehicle_id/history", locationHandler.GetLocationHistory)
	}

	return r
}
