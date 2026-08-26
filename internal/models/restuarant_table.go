package models

import (
	"time"

	"github.com/google/uuid"
)

type RestaurantTable struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RestaurantID uuid.UUID `gorm:"type:uuid;not null;index"`

	TableNumber int `gorm:"not null"`
	Capacity    int    `gorm:"not null"`
IsAvailable bool `gorm:"default:true"`
	

	CreatedAt time.Time
	UpdatedAt time.Time
}