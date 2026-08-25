package models


import "time"

type UserRole string

const (
    RoleCustomer UserRole = "customer"
    RoleStaff    UserRole = "staff"
    RoleAdmin    UserRole = "admin"
)

type User struct {
    ID           int64
    Name         string
    Email        string
    PasswordHash string
    Role         UserRole
    RestaurantID *int64
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type RefreshToken struct {
    ID        int64
    UserID    int64
    TokenHash string
    ExpiresAt time.Time
    Revoked   bool
    CreatedAt time.Time
}