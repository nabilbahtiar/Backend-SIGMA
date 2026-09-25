package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SensorConfiguration struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SensorID      uuid.UUID `gorm:"type:uuid;not null;index" json:"sensor_id"`
	MinThreshold  float64   `json:"min_threshold"`
	MaxThreshold  float64   `json:"max_threshold"`
	SeverityLevel string    `gorm:"not null" json:"severity_level"` // WARNING, CRITICAL, EMERGENCY
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	UpdatedBy     uuid.UUID `gorm:"type:uuid" json:"updated_by"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (s *SensorConfiguration) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
