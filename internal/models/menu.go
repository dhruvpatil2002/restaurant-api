package models

import (
	"time"

	"github.com/google/uuid"
)

type Menu struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	RestaurantID uuid.UUID `gorm:"type:uuid;not null;index"`

	Name        string  `gorm:"not null"`
	Category    string  `gorm:"not null;index"`
	Description string
	Price       float64 `gorm:"not null"`

	Image string

	IsAvailable bool `gorm:"default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}