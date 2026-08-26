package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"restaurant-backend/internal/models"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		DB: db,
	}
}

// =====================================================
// CREATE ORDER
// =====================================================

func (r *OrderRepository) Create(
	order *models.Order,
) error {

	return r.DB.Create(order).Error
}

// =====================================================
// GET ORDER BY ID
// =====================================================

func (r *OrderRepository) FindByID(
	id uuid.UUID,
) (*models.Order, error) {

	var order models.Order

	err := r.DB.
		Preload("Items").
		Where("id = ?", id).
		First(&order).
		Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

// =====================================================
// GET USER ORDERS
// =====================================================

func (r *OrderRepository) FindByUserID(
	userID uuid.UUID,
) ([]models.Order, error) {

	var orders []models.Order

	err := r.DB.
		Preload("Items").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).
		Error

	return orders, err
}

// =====================================================
// GET RESTAURANT ORDERS
// =====================================================

func (r *OrderRepository) FindByRestaurantID(
	restaurantID uuid.UUID,
) ([]models.Order, error) {

	var orders []models.Order

	err := r.DB.
		Preload("Items").
		Where("restaurant_id = ?", restaurantID).
		Order("created_at DESC").
		Find(&orders).
		Error

	return orders, err
}

// =====================================================
// UPDATE
// =====================================================

func (r *OrderRepository) Update(
	order *models.Order,
) error {

	return r.DB.Save(order).Error
}

// =====================================================
// DELETE
// =====================================================

func (r *OrderRepository) Delete(
	id uuid.UUID,
) error {

	return r.DB.Delete(
		&models.Order{},
		"id = ?",
		id,
	).Error
}