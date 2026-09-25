package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sensor struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SensorCode string    `gorm:"uniqueIndex;not null" json:"sensor_code"`
	SensorName string    `gorm:"not null" json:"sensor_name"`
	SensorType string    `gorm:"not null" json:"sensor_type"`
	Unit       string    `json:"unit"` // e.g., Celcius, %, Boolean
	Location   string    `json:"location"`
	Status     string    `gorm:"default:'ACTIVE'" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (s *Sensor) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
