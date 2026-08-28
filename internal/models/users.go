package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleOwner    UserRole = "owner"
	RoleStaff    UserRole = "staff"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string     `gorm:"not null"`
	Email        string     `gorm:"uniqueIndex;not null"`
	PasswordHash string     `gorm:"not null"`
	Role UserRole `gorm:"type:varchar(20);not null;default:'customer'" json:"role"`
	RestaurantID *uuid.UUID `gorm:"type:uuid;index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RefreshToken struct {
	ID        int64     `gorm:"primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Revoked   bool      `gorm:"not null;default:false"`
	CreatedAt time.Time
}