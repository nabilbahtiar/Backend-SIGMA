package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    *uuid.UUID `gorm:"type:uuid" json:"user_id"` // Opsional jika broadcast ke grup
	AlarmID   uuid.UUID  `gorm:"type:uuid;not null" json:"alarm_id"`
	Message   string     `gorm:"not null" json:"message"`
	Status    string     `gorm:"default:'PENDING'" json:"status"` // PENDING, SENT, FAILED
	SentAt    *time.Time `json:"sent_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return
}
