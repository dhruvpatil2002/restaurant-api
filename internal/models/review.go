package models

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RestaurantID uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index"`

	Rating  int    `gorm:"not null"`
	Comment string

	CreatedAt time.Time
	UpdatedAt time.Time
}