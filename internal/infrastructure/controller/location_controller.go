package controller

import (
	"FleetManagementSystem/internal/entity"
	"FleetManagementSystem/internal/infrastructure/api/location"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LocationController struct {
	service location.Service
}

func NewLocationController(service location.Service) *LocationController {
	return &LocationController{service: service}

}

func (h *LocationController) GetLatestLocation(c *gin.Context) {

	// vehicleID := c.Param("vehicle_id")
	// TODO: Implement logic to get latest location from repository
	c.JSON(http.StatusOK, entity.ConvertToBaseResponse("success", http.StatusOK, entity.EmptyObj{}))
}
