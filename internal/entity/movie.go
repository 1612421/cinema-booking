package entity

import (
	"github.com/google/uuid"
	"time"
)

type Movie struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Duration  int       `json:"duration,omitempty"` // minutes
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Movie) TableName() string {
	return "movies"
}
