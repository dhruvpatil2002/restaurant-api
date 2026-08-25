package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RestaurantID uuid.UUID `gorm:"type:uuid;not null;index"`

	Name        string `gorm:"not null"`
	Description string

	CreatedAt time.Time
	UpdatedAt time.Time
}