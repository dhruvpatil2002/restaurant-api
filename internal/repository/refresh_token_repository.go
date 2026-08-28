package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"restaurant-backend/internal/models"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, rt *models.RefreshToken) error {
	return r.db.WithContext(ctx).Create(rt).Error
}

func (r *RefreshTokenRepository) GetByHashAndNotRevoked(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND revoked = ?", tokenHash, false).
		First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &rt, nil
}

func (r *RefreshTokenRepository) GetByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &rt, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, rt *models.RefreshToken) error {
	rt.Revoked = true
	return r.db.WithContext(ctx).Save(rt).Error
}