package models

import (
    "time"

    "github.com/google/uuid"
)

type UserRole string

const (
    RoleCustomer UserRole = "customer"
    RoleStaff    UserRole = "staff"
    RoleAdmin    UserRole = "admin"
)

type User struct {
    ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Name         string
    Email        string
    PasswordHash string   `gorm:"not null"`
    Role string `gorm:"type:varchar(20);not null;default:'customer'"`
    RestaurantID *uuid.UUID // better to keep this uuid too; if you don't use it yet, can leave as is for now
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type RefreshToken struct {
    ID        int64     `gorm:"primaryKey"`
    UserID    uuid.UUID `gorm:"type:uuid;not null;index"` // change from int64 to uuid.UUID
    TokenHash string    `gorm:"uniqueIndex;not null"`
    ExpiresAt time.Time `gorm:"not null"`
    Revoked   bool      `gorm:"not null;default:false"`
    CreatedAt time.Time
}