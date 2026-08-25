package models

import (
	"time"

	"github.com/google/uuid"
)

type MenuItem struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RestaurantID uuid.UUID `gorm:"type:uuid;not null;index"`
	CategoryID   uuid.UUID `gorm:"type:uuid;index"`

	Name        string  `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`

	Image string

	IsAvailable bool `gorm:"default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}