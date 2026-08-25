package models

import (
	"time"

	"github.com/google/uuid"
)

type Restaurant struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null;index"`

	Name        string `gorm:"not null"`
	Description string

	Address string `gorm:"not null"`
	City    string `gorm:"not null"`
	State   string `gorm:"not null"`
	Pincode string `gorm:"not null"`

	Phone string
	Email string
	Image string

	OpeningTime string
	ClosingTime string

	IsOpen bool `gorm:"default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}