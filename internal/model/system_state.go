package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SystemState struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CurrentMode string    `gorm:"not null" json:"current_mode"` // ARMED, MAINTENANCE
	SetBy       uuid.UUID `gorm:"type:uuid;not null" json:"set_by"`
	Reason      string    `json:"reason"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (s *SystemState) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
