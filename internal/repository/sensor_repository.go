package repository

import (
	"server-room-auth/internal/model"

	"gorm.io/gorm"
)

type SensorRepository struct {
	db *gorm.DB
}

func NewSensorRepository(db *gorm.DB) *SensorRepository {
	return &SensorRepository{db: db}
}

// GetAllSensorsWithConfig mengambil semua sensor beserta konfigurasi ambang batasnya
func (r *SensorRepository) GetAllSensorsWithConfig() ([]model.Sensor, error) {
	var sensors []model.Sensor
	// Karena kita belum mengatur foreign key HasMany di model Sensor ke Configuration,
	// kita bisa mengambilnya secara terpisah atau menggunakan preload jika sudah di-set.
	// Untuk saat ini, kita ambil semua sensor saja.
	err := r.db.Find(&sensors).Error
	return sensors, err
}

// GetConfigurationsBySensorID mengambil konfigurasi aktif untuk sebuah sensor
func (r *SensorRepository) GetConfigurationsBySensorID(sensorID string) ([]model.SensorConfiguration, error) {
	var configs []model.SensorConfiguration
	err := r.db.Where("sensor_id = ? AND is_active = ?", sensorID, true).Find(&configs).Error
	return configs, err
}

// GetLatestTelemetry mengambil 1 data telemetri terbaru dari sebuah sensor
func (r *SensorRepository) GetLatestTelemetry(sensorID string) (*model.SensorTelemetry, error) {
	var telemetry model.SensorTelemetry
	err := r.db.Where("sensor_id = ?", sensorID).Order("recorded_at desc").First(&telemetry).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Belum ada data
		}
		return nil, err
	}
	return &telemetry, nil
}
