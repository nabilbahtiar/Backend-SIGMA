package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HelpdeskTicket struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	TicketNumber    string     `gorm:"uniqueIndex;not null" json:"ticket_number"`
	AlarmID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"alarm_id"`
	AssignedBy      *uuid.UUID `gorm:"type:uuid" json:"assigned_by"`
	AssignedTo      *uuid.UUID `gorm:"type:uuid" json:"assigned_to"`
	Title           string     `gorm:"not null" json:"title"`
	Description     string     `json:"description"`
	Status          string     `gorm:"default:'OPEN'" json:"status"` // OPEN, IN_PROGRESS, RESOLVED, CLOSED
	ImageBefore     string     `json:"image_before"`
	ImageAfter      string     `json:"image_after"`
	ResolutionNotes string     `json:"resolution_notes"`
	ClosedBy        *uuid.UUID `gorm:"type:uuid" json:"closed_by"`
	AssignedAt      *time.Time `json:"assigned_at"`
	ClosedAt        *time.Time `json:"closed_at"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (h *HelpdeskTicket) BeforeCreate(tx *gorm.DB) (err error) {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return
}
