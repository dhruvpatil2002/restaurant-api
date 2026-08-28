package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"restaurant-backend/internal/models"
)

type RestaurantRepository struct {
	DB *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) *RestaurantRepository {
	return &RestaurantRepository{
		DB: db,
	}
}

// =====================================================
// CREATE
// =====================================================

func (r *RestaurantRepository) Create(
	restaurant *models.Restaurant,
) error {

	return r.DB.Create(restaurant).Error
}

// =====================================================
// FIND BY ID
// =====================================================

func (r *RestaurantRepository) FindByID(
	id uuid.UUID,
) (*models.Restaurant, error) {

	var restaurant models.Restaurant

	err := r.DB.
		Where("id = ?", id).
		First(&restaurant).Error

	if err != nil {
		return nil, err
	}

	return &restaurant, nil
}

// =====================================================
// FIND BY OWNER
// =====================================================

func (r *RestaurantRepository) FindByOwnerID(
	ownerID uuid.UUID,
) (*models.Restaurant, error) {

	var restaurant models.Restaurant

	err := r.DB.
		Where("owner_id = ?", ownerID).
		First(&restaurant).Error

	if err != nil {
		return nil, err
	}

	return &restaurant, nil
}

// =====================================================
// FIND ALL
// =====================================================

func (r *RestaurantRepository) FindAll() (
	[]models.Restaurant,
	error,
) {

	var restaurants []models.Restaurant

	err := r.DB.
		Order("created_at DESC").
		Find(&restaurants).Error

	return restaurants, err
}

// =====================================================
// UPDATE
// =====================================================

func (r *RestaurantRepository) Update(
	restaurant *models.Restaurant,
) error {

	return r.DB.Save(restaurant).Error
}

// =====================================================
// DELETE
// =====================================================

func (r *RestaurantRepository) Delete(
	id uuid.UUID,
) error {

	return r.DB.
		Delete(
			&models.Restaurant{},
			"id = ?",
			id,
		).Error
}

// =====================================================
// PAGINATION
// =====================================================

func (r *RestaurantRepository) FindAllPaginated(
	page int,
	limit int,
	search string,
) ([]models.Restaurant, int64, error) {

	var restaurants []models.Restaurant
	var total int64

	offset := (page - 1) * limit

	query := r.DB.Model(&models.Restaurant{})

	// Search restaurant name
	if search != "" {
		query = query.Where(
			"name ILIKE ?",
			"%"+search+"%",
		)
	}

	// Total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Data
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&restaurants).Error; err != nil {

		return nil, 0, err
	}

	return restaurants, total, nil
}