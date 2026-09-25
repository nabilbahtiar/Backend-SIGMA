package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User merepresentasikan tabel users di database
type User struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	NIK            string    `gorm:"uniqueIndex;not null" json:"nik"`
	Nama           string    `gorm:"not null" json:"nama"`
	Jabatan        string    `gorm:"default:''" json:"jabatan"`
	Unit           string    `gorm:"default:''" json:"unit"`
	TipePegawai    string    `gorm:"default:''" json:"tipe_pegawai"`
	NoHP           string    `gorm:"default:''" json:"no_hp"`
	IDTelegram     string    `gorm:"default:''" json:"id_telegram"`
	TelegramChatID string    `gorm:"default:''" json:"telegram_chat_id"` // ID numerik untuk bot mengirim pesan
	PasswordHash   string    `gorm:"not null" json:"-"`
	Role           string    `gorm:"not null" json:"role"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// GetNIK mengembalikan NIK
func (u *User) GetNIK() string {
	return u.NIK
}

// BeforeCreate adalah hook GORM untuk meng-generate UUID otomatis sebelum disimpan
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
