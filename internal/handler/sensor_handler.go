package handler

import (
	"net/http"

	"server-room-auth/internal/service"

	"github.com/gin-gonic/gin"
)

type SensorHandler struct {
	sensorService *service.SensorService
}

func NewSensorHandler(service *service.SensorService) *SensorHandler {
	return &SensorHandler{sensorService: service}
}

// GetDashboardSensors mengembalikan data 6 sensor (nilai terbaru + ambang batas)
func (h *SensorHandler) GetDashboardSensors(c *gin.Context) {
	data, err := h.sensorService.GetDashboardData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data sensor dashboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data dashboard",
		"data":    data,
	})
}
