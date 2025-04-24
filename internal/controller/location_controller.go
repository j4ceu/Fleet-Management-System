package controller

import (
	"FleetManagementSystem/internal/api/location"
	"FleetManagementSystem/internal/entity"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type LocationController struct {
	service location.Service
}

func NewLocationController(service location.Service) *LocationController {
	return &LocationController{service: service}

}

func (h *LocationController) GetLatestLocation(c *gin.Context) {

	vehicleID := c.Param("vehicle_id")

	location, err := h.service.GetLatestLocation(c, vehicleID)
	if err != nil {
		baseResponse := entity.ConvertErrorToBaseResponse("failed", http.StatusInternalServerError, entity.EmptyObj{}, err.Error())
		c.JSON(http.StatusInternalServerError, baseResponse)
		return
	}
	c.JSON(http.StatusOK, entity.ConvertToBaseResponse("success", http.StatusOK, location))
}

func (h *LocationController) GetLocationHistory(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	startStr := c.DefaultQuery("start", "0")
	endStr := c.DefaultQuery("end", strconv.FormatInt(time.Now().Unix(), 10))

	start, err1 := strconv.ParseInt(startStr, 10, 64)
	end, err2 := strconv.ParseInt(endStr, 10, 64)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, entity.ConvertErrorToBaseResponse("Invalid time range", http.StatusBadRequest, nil, err1.Error()+" | "+err2.Error()))
		return
	}

	payload := entity.LocationHistoryPayload{
		VehicleID: vehicleID,
		StartTime: start,
		EndTime:   end,
	}

	locations, err := h.service.GetLocationHistory(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ConvertErrorToBaseResponse("Failed to get history", http.StatusInternalServerError, nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, entity.ConvertToBaseResponse("OK", http.StatusOK, locations))
}
