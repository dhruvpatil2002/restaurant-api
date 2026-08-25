package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RestaurantID uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index"`

	TotalAmount float64 `gorm:"not null"`

	Status string `gorm:"default:'pending'"`

	Notes string

	CreatedAt time.Time
	UpdatedAt time.Time

	Items []OrderItem `gorm:"foreignKey:OrderID"`
}