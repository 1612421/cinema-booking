package entity

import (
	"github.com/google/uuid"
	"time"
)

type Seat struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id,omitempty"`
	ScreenID  uuid.UUID `json:"screen_id,omitempty"` // minutes
	Row       string    `gorm:"column:seat_row" json:"row,omitempty"`
	Number    int       `json:"number,omitempty"`
	Class     string    `json:"class,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SeatWithStatus struct {
	*Seat
	IsAvailable bool `json:"is_available"`
}

func (s *Seat) TableName() string {
	return "seats"
}
