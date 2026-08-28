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

func (r *MenuRepository) Update(
	menu *models.Menu,
) error {
	return r.DB.Save(menu).Error
}


func (r *MenuRepository) Delete(
	id uuid.UUID,
) error {
	return r.DB.Delete(
		&models.Menu{},
		"id = ?",
		id,
	).Error
}

func (r *MenuRepository) FindPaginated(
	restaurantID uuid.UUID,
	page int,
	limit int,
	search string,
	category string,
	available *bool,
) ([]models.Menu, int64, error) {

	var menus []models.Menu
	var total int64

	offset := (page - 1) * limit

	query := r.DB.
		Model(&models.Menu{}).
		Where(
			"restaurant_id = ?",
			restaurantID,
		)


	if search != "" {
		query = query.Where(
			"(name ILIKE ? OR description ILIKE ?)",
			"%"+search+"%",
			"%"+search+"%",
		)
	}


	if category != "" {
		query = query.Where(
			"category = ?",
			category,
		)
	}
	
	if available != nil {
		query = query.Where(
			"is_available = ?",
			*available,
		)
	}

	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&menus).Error; err != nil {

		return nil, 0, err
	}

	return menus, total, nil
}