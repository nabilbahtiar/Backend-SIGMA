package service

import (
	"server-room-auth/internal/repository"
	"time"
)

type SensorService struct {
	sensorRepo *repository.SensorRepository
}

func NewSensorService(repo *repository.SensorRepository) *SensorService {
	return &SensorService{sensorRepo: repo}
}

// DashboardSensorResponse adalah struktur data (JSON) yang akan dikirim ke Frontend
type DashboardSensorResponse struct {
	SensorID       string    `json:"sensor_id"`
	SensorCode     string    `json:"sensor_code"`
	SensorName     string    `json:"sensor_name"`
	Unit           string    `json:"unit"`
	CurrentValue   float64   `json:"current_value"`
	Status         string    `json:"status"` // NORMAL, WARNING, CRITICAL, EMERGENCY
	MinThreshold   float64   `json:"min_threshold"`
	MaxThreshold   float64   `json:"max_threshold"`
	LastUpdated    time.Time `json:"last_updated"`
}

func (s *SensorService) GetDashboardData() ([]DashboardSensorResponse, error) {
	// 1. Ambil semua sensor
	sensors, err := s.sensorRepo.GetAllSensorsWithConfig()
	if err != nil {
		return nil, err
	}

	var response []DashboardSensorResponse

	// 2. Loop setiap sensor untuk menyusun data dashboard
	for _, sensor := range sensors {
		sensorIDStr := sensor.ID.String()

		// Ambil telemetry terbaru
		telemetry, _ := s.sensorRepo.GetLatestTelemetry(sensorIDStr)
		currentVal := 0.0
		lastUpdate := sensor.UpdatedAt
		if telemetry != nil {
			currentVal = telemetry.SensorValue
			lastUpdate = telemetry.RecordedAt
		}

		// Ambil konfigurasi (threshold)
		configs, _ := s.sensorRepo.GetConfigurationsBySensorID(sensorIDStr)
		
		// Secara default, kita asumsikan NORMAL
		status := "NORMAL"
		var activeMin, activeMax float64

		// Evaluasi logika status berdasarkan threshold
		for _, config := range configs {
			// Simpan nilai threshold untuk ditampilkan di frontend (bisa diperbaiki logikanya agar lebih fleksibel)
			if config.SeverityLevel == "WARNING" {
				activeMin = config.MinThreshold
				activeMax = config.MaxThreshold
			}

			// Cek apakah nilai saat ini melanggar threshold
			if currentVal < config.MinThreshold || currentVal > config.MaxThreshold {
				// Timpa status dengan level bahaya tertinggi (asumsi loop menaikkan bahaya)
				status = config.SeverityLevel
			}
		}

		// 3. Susun ke dalam format respons
		card := DashboardSensorResponse{
			SensorID:     sensorIDStr,
			SensorCode:   sensor.SensorCode,
			SensorName:   sensor.SensorName,
			Unit:         sensor.Unit,
			CurrentValue: currentVal,
			Status:       status,
			MinThreshold: activeMin,
			MaxThreshold: activeMax,
			LastUpdated:  lastUpdate,
		}
		response = append(response, card)
	}

	return response, nil
}
