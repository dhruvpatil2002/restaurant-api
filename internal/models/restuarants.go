package models

import (
	"time"

	"github.com/google/uuid"
)

type Restaurant struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OwnerID uuid.UUID `json:"owner_id" gorm:"type:uuid;not null;index"`

	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`

	Address string `json:"address" gorm:"not null"`
	City    string `json:"city" gorm:"not null"`
	State   string `json:"state" gorm:"not null"`
	Pincode string `json:"pincode" gorm:"not null"`

	Phone string `json:"phone"`
	Email string `json:"email"`
	Image string `json:"image"`

	OpeningTime string `json:"opening_time"`
	ClosingTime string `json:"closing_time"`

	IsOpen bool `json:"is_open" gorm:"default:true"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}