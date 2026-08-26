package models

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	RestaurantID uuid.UUID `gorm:"type:uuid;not null;index"`

	TableID uuid.UUID `gorm:"type:uuid;not null;index"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`

	ReservationTime time.Time `gorm:"not null;index"`

	GuestCount int `gorm:"not null"`

	Status string `gorm:"type:varchar(20);not null;default:'pending';index"`

	Notes string

	CreatedAt time.Time
	UpdatedAt time.Time
}