package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Alarm struct {
	ID                uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	SensorID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"sensor_id"`
	SensorTelemetryID uuid.UUID  `gorm:"type:uuid;not null" json:"sensor_telemetry_id"` // Referensi ke data yang memicu alarm
	AlarmLevel        string     `gorm:"not null" json:"alarm_level"`                   // WARNING, CRITICAL, EMERGENCY
	Status            string     `gorm:"default:'ACTIVE'" json:"status"`                // ACTIVE, ACKNOWLEDGED, RESOLVED
	TriggeredAt       time.Time  `gorm:"autoCreateTime" json:"triggered_at"`
	ResolvedAt        *time.Time `json:"resolved_at,omitempty"`
}

func (a *Alarm) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
