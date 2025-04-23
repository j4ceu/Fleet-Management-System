package config

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// apiV1 := r.Group("/api/v1")

	// vehicles := apiV1.Group("/vehicles")

	{
		// vehicles.GET("/",)
	}

	return r
}
