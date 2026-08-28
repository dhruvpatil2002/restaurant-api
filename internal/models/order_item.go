package models

import (
	"github.com/google/uuid"
)

type OrderItem struct {
    ID         uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    OrderID    uuid.UUID `json:"order_id" gorm:"type:uuid;not null;index"`
    MenuItemID uuid.UUID `json:"menu_item_id" gorm:"type:uuid;not null"`
    Quantity   int       `json:"quantity" gorm:"not null"`
    Price      float64   `json:"price" gorm:"not null"`
}