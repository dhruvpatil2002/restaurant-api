package models

import (
	"time"

	"github.com/google/uuid"
)

type RestaurantTable struct {
    ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    RestaurantID uuid.UUID `gorm:"type:uuid;not null;index" json:"restaurant_id"`

    TableNumber  int  `gorm:"not null" json:"table_number"`
    Capacity     int  `gorm:"not null" json:"capacity"`
    IsAvailable  bool `gorm:"default:true" json:"is_available"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}