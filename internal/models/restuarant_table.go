package models

import (
	"time"

	"github.com/google/uuid"
)

type RestaurantTable struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RestaurantID uuid.UUID `gorm:"type:uuid;not null;index"`

	TableNumber string `gorm:"not null"`
	Capacity    int    `gorm:"not null"`

	Status string `gorm:"default:'available'"`

	CreatedAt time.Time
	UpdatedAt time.Time
}