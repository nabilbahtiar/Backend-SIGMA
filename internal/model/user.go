package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User merepresentasikan tabel users di database
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	NIK          string    `gorm:"uniqueIndex;not null" json:"nik"`
	Nama         string    `gorm:"not null" json:"nama"`
	Jabatan      string    `gorm:"not null" json:"jabatan"`
	Unit         string    `gorm:"not null" json:"unit"`
	TipePegawai  string    `gorm:"not null" json:"tipe_pegawai"`
	NoHP         string    `gorm:"not null" json:"no_hp"`
	PasswordHash string    `gorm:"not null" json:"-"` // Tidak dikirim ke klien
	Role         string    `gorm:"not null" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BeforeCreate adalah hook GORM untuk meng-generate UUID otomatis sebelum disimpan
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
