package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      *uuid.UUID `gorm:"type:uuid" json:"user_id"` // Opsional, bisa null jika dari sistem otomatis
	Activity    string     `gorm:"not null" json:"activity"`
	Module      string     `json:"module"`
	Description string     `json:"description"`
	IPAddress   string     `json:"ip_address"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
