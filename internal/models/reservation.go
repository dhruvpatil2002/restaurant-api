package models

import (
	"time"

	"github.com/google/uuid"
)



type Reservation struct {
    ID             uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    RestaurantID   uuid.UUID `json:"restaurant_id" gorm:"type:uuid;not null;index"`
    TableID        uuid.UUID `json:"table_id" gorm:"type:uuid;not null;index"`
    UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
    ReservationTime time.Time `json:"reservation_time" gorm:"column:reservation_time;not null;index"`
    GuestCount     int       `json:"guest_count" gorm:"not null"`
    Status         string    `json:"status" gorm:"type:varchar(20);not null;default:'pending';index"`
    Notes          string    `json:"notes"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}