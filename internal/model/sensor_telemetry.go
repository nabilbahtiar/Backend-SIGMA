package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SensorTelemetry struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SensorID    uuid.UUID `gorm:"type:uuid;not null;index" json:"sensor_id"`
	SensorValue float64   `gorm:"not null" json:"sensor_value"`
	RecordedAt  time.Time `gorm:"autoCreateTime" json:"recorded_at"`
}

func (s *SensorTelemetry) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
