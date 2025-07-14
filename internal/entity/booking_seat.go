package entity

import (
	"github.com/google/uuid"
	"time"
)

type BookingSeat struct {
	ID         uuid.UUID `gorm:"primaryKey;not null" json:"id,omitempty"`
	BookingID  uuid.UUID `json:"booking_id,omitempty"`
	ShowtimeID uuid.UUID `json:"showtime_id,omitempty"`
	SeatID     uuid.UUID `json:"seat_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (b *BookingSeat) TableName() string {
	return "booking_seats"
}
