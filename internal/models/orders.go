package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
    ID           uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    RestaurantID uuid.UUID `json:"restaurant_id" gorm:"type:uuid;not null;index"`
    UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
    TotalAmount  float64   `json:"total_amount" gorm:"not null"`
    Status       string    `json:"status" gorm:"default:'pending'"`
    Notes        string    `json:"notes"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    Items        []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}