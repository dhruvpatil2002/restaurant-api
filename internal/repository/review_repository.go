package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"restaurant-backend/internal/models"
)

type ReviewRepository struct {
	DB *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{
		DB: db,
	}
}


func (r *ReviewRepository) Create(
	review *models.Review,
) error {

	return r.DB.Create(review).Error
}


func (r *ReviewRepository) FindByID(
	id uuid.UUID,
) (*models.Review, error) {

	var review models.Review

	err := r.DB.
		Where("id = ?", id).
		First(&review).
		Error

	if err != nil {
		return nil, err
	}

	return &review, nil
}




func (r *ReviewRepository) FindByRestaurantID(
	restaurantID uuid.UUID,
) ([]models.Review, error) {

	var reviews []models.Review

	err := r.DB.
		Where("restaurant_id = ?", restaurantID).
		Order("created_at DESC").
		Find(&reviews).
		Error

	return reviews, err
}



func (r *ReviewRepository) FindByUserID(
	userID uuid.UUID,
) ([]models.Review, error) {

	var reviews []models.Review

	err := r.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reviews).
		Error

	return reviews, err
}



func (r *ReviewRepository) FindByUserAndRestaurant(
	userID uuid.UUID,
	restaurantID uuid.UUID,
) (*models.Review, error) {

	var review models.Review

	err := r.DB.
		Where(
			"user_id = ? AND restaurant_id = ?",
			userID,
			restaurantID,
		).
		First(&review).
		Error

	if err != nil {
		return nil, err
	}

	return &review, nil
}



func (r *ReviewRepository) Update(
	review *models.Review,
) error {

	return r.DB.Save(review).Error
}


func (r *ReviewRepository) Delete(
	id uuid.UUID,
) error {

	return r.DB.Delete(
		&models.Review{},
		"id = ?",
		id,
	).Error
}