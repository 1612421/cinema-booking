package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID `gorm:"primaryKey;not null" json:"id,omitempty"`
	Username    string    `gorm:"uniqueIndex" json:"username,omitempty"`
	Password    string    `json:"-"`
	Status      string    `json:"status,omitempty"` // active,ban,disabled
	Address     string    `json:"address,omitempty"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u *User) TableName() string {
	return "users"
}
