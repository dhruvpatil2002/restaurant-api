package repository

import (
    "restaurant-backend/internal/models"
    "github.com/google/uuid"

    "gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}


type RefreshTokenRepository struct {
    db *gorm.DB
}
func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
    var u models.User
    if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
    var u models.User
    if err := r.db.First(&u, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &u, nil
}
func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
    return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(rt *models.RefreshToken) error {
    return r.db.Create(rt).Error
}

func (r *RefreshTokenRepository) GetByHashAndNotRevoked(tokenHash string) (*models.RefreshToken, error) {
    var rt models.RefreshToken
    if err := r.db.Where("token_hash = ? AND revoked = ?", tokenHash, false).First(&rt).Error; err != nil {
        return nil, err
    }
    return &rt, nil
}

func (r *RefreshTokenRepository) GetByHash(tokenHash string) (*models.RefreshToken, error) {
    var rt models.RefreshToken
    if err := r.db.Where("token_hash = ?", tokenHash).First(&rt).Error; err != nil {
        return nil, err
    }
    return &rt, nil
}

func (r *RefreshTokenRepository) Revoke(rt *models.RefreshToken) error {
    rt.Revoked = true
    return r.db.Save(rt).Error
}