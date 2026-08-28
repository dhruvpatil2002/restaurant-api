package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"restaurant-backend/internal/models"
)

type ReservationRepository struct {
	DB *gorm.DB
}

func NewReservationRepository(
	db *gorm.DB,
) *ReservationRepository {

	return &ReservationRepository{
		DB: db,
	}
}

// =====================================================
// CREATE
// =====================================================

func (r *ReservationRepository) Create(
	reservation *models.Reservation,
) error {

	return r.DB.Create(reservation).Error
}

// =====================================================
// FIND BY ID
// =====================================================

func (r *ReservationRepository) FindByID(
	id uuid.UUID,
) (*models.Reservation, error) {

	var reservation models.Reservation

	err := r.DB.
		Where("id = ?", id).
		First(&reservation).
		Error

	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

// =====================================================
// GET USER RESERVATIONS
// =====================================================

func (r *ReservationRepository) FindByUserID(
	userID uuid.UUID,
) ([]models.Reservation, error) {

	var reservations []models.Reservation

	err := r.DB.
		Where("user_id = ?", userID).
		Order("reservation_time DESC").
		Find(&reservations).
		Error

	return reservations, err
}

// =====================================================
// GET RESTAURANT RESERVATIONS
// =====================================================

func (r *ReservationRepository) FindByRestaurantID(
	restaurantID uuid.UUID,
) ([]models.Reservation, error) {

	var reservations []models.Reservation

	err := r.DB.
		Where("restaurant_id = ?", restaurantID).
		Order("reservation_time DESC").
		Find(&reservations).
		Error

	return reservations, err
}

// =====================================================
// CHECK TABLE CONFLICT
// =====================================================

func (r *ReservationRepository) HasConflict(
	tableID uuid.UUID,
	reservationTime time.Time,
) (bool, error) {

	var count int64

	err := r.DB.
		Model(&models.Reservation{}).
		Where(
			`
			table_id = ?
			AND reservation_time = ?
			AND status IN ?
			`,
			tableID,
			reservationTime,
			[]string{"pending", "confirmed"},
		).
		Count(&count).
		Error

	return count > 0, err
}

// =====================================================
// UPDATE
// =====================================================

func (r *ReservationRepository) Update(
	reservation *models.Reservation,
) error {

	return r.DB.Save(reservation).Error
}

// =====================================================
// DELETE
// =====================================================

func (r *ReservationRepository) Delete(
	id uuid.UUID,
) error {

	return r.DB.
		Delete(
			&models.Reservation{},
			"id = ?",
			id,
		).
		Error
}