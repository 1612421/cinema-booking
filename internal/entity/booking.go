package entity

import (
	"github.com/google/uuid"
	"time"
)

const (
	BookingStatusPending   = "pending"
	BookingStatusConfirmed = "confirmed"
	BookingStatusCanceled  = "canceled"
)

type Booking struct {
	ID         uuid.UUID `gorm:"primaryKey;not null" json:"id,omitempty"`
	ShowtimeID uuid.UUID `json:"showtime_id,omitempty"`
	UserId     uuid.UUID `json:"user_id,omitempty"` // minutes
	Status     string    `json:"status,omitempty"`  // pending,confirmed,canceled
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (b *Booking) TableName() string {
	return "bookings"
}
