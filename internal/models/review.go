package models

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
    ID           uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    RestaurantID uuid.UUID `json:"restaurant_id" gorm:"type:uuid;not null;index"`
    UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
    Rating       int       `json:"rating" gorm:"not null"`
    Comment      string    `json:"comment"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}