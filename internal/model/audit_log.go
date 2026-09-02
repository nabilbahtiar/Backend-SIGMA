package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLog merepresentasikan tabel audit_logs di database untuk melacak aktivitas user
type AuditLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Timestamp  time.Time `gorm:"index;not null" json:"timestamp"` // Diindeks agar pencarian log berdasarkan waktu cepat
	Username   string    `gorm:"index;default:'Guest'" json:"username"` // 'Guest' jika belum login
	Action     string    `gorm:"not null" json:"action"`          // Contoh: "API Access", "Login Failed"
	Method     string    `gorm:"not null" json:"method"`
	Path       string    `gorm:"not null" json:"path"`
	StatusCode int       `gorm:"not null" json:"status_code"`
	ClientIP   string    `gorm:"not null" json:"client_ip"`
	UserAgent  string    `gorm:"not null" json:"user_agent"`      // Aplikasi/Browser yang digunakan
	Latency    string    `gorm:"not null" json:"latency"`
}

// BeforeCreate adalah hook GORM untuk meng-generate UUID otomatis sebelum disimpan
func (a *AuditLog) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
