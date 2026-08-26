package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"restaurant-backend/internal/models"
)

type MenuRepository struct {
	DB *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{
		DB: db,
	}
}

// Create menu item
func (r *MenuRepository) Create(
	menu *models.Menu,
) error {
	return r.DB.Create(menu).Error
}

// Get menu item by ID
func (r *MenuRepository) FindByID(
	id uuid.UUID,
) (*models.Menu, error) {

	var menu models.Menu

	err := r.DB.
		Where("id = ?", id).
		First(&menu).Error

	if err != nil {
		return nil, err
	}

	return &menu, nil
}

// Get all menu items of a restaurant
func (r *MenuRepository) FindByRestaurantID(
	restaurantID uuid.UUID,
) ([]models.Menu, error) {

	var menus []models.Menu

	err := r.DB.
		Where("restaurant_id = ?", restaurantID).
		Order("created_at DESC").
		Find(&menus).Error

	return menus, err
}

// Get menu by category
func (r *MenuRepository) FindByCategory(
	restaurantID uuid.UUID,
	category string,
) ([]models.Menu, error) {

	var menus []models.Menu

	err := r.DB.
		Where(
			"restaurant_id = ? AND category = ?",
			restaurantID,
			category,
		).
		Order("created_at DESC").
		Find(&menus).Error

	return menus, err
}

// Update
func (r *MenuRepository) Update(
	menu *models.Menu,
) error {
	return r.DB.Save(menu).Error
}

// Delete
func (r *MenuRepository) Delete(
	id uuid.UUID,
) error {
	return r.DB.Delete(
		&models.Menu{},
		"id = ?",
		id,
	).Error
}