package models

import (
	"github.com/google/uuid"
)

type OrderItem struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	OrderID    uuid.UUID `gorm:"type:uuid;not null;index"`
	MenuItemID uuid.UUID `gorm:"type:uuid;not null"`

	Quantity int `gorm:"not null"`

	Price float64 `gorm:"not null"`
}