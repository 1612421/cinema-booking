package entity

import (
	"github.com/google/uuid"
	"time"
)

type Showtime struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id,omitempty"`
	MovieID   uuid.UUID `json:"movie_id,omitempty"`
	ScreenID  uuid.UUID `json:"screen_id,omitempty"` // minutes
	StartTime time.Time `json:"start_time,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Showtime) TableName() string {
	return "showtime"
}
