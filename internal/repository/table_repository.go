package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"restaurant-backend/internal/models"
)

type TableRepository struct {
	DB *gorm.DB
}

func NewTableRepository(db *gorm.DB) *TableRepository {
	return &TableRepository{
		DB: db,
	}
}

// Create table
func (r *TableRepository) Create(
	table *models.RestaurantTable,
) error {

	return r.DB.Create(table).Error
}

// Find table by ID
func (r *TableRepository) FindByID(
	id uuid.UUID,
) (*models.RestaurantTable, error) {

	var table models.RestaurantTable

	err := r.DB.
		Where("id = ?", id).
		First(&table).
		Error

	if err != nil {
		return nil, err
	}

	return &table, nil
}

// Get all tables of restaurant
func (r *TableRepository) FindByRestaurantID(
	restaurantID uuid.UUID,
) ([]models.RestaurantTable, error) {

	var tables []models.RestaurantTable

	err := r.DB.
		Where("restaurant_id = ?", restaurantID).
		Order("table_number ASC").
		Find(&tables).
		Error

	return tables, err
}

// Update table
func (r *TableRepository) Update(
	table *models.RestaurantTable,
) error {

	return r.DB.Save(table).Error
}

// Delete table
func (r *TableRepository) Delete(
	id uuid.UUID,
) error {

	return r.DB.
		Delete(
			&models.RestaurantTable{},
			"id = ?",
			id,
		).
		Error
}