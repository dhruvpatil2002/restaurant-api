package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"
)

type ReservationService struct {
	ReservationRepo *repository.ReservationRepository
	RestaurantRepo  *repository.RestaurantRepository
	TableRepo       *repository.TableRepository
}

func NewReservationService(
	reservationRepo *repository.ReservationRepository,
	restaurantRepo *repository.RestaurantRepository,
	tableRepo *repository.TableRepository,
) *ReservationService {

	return &ReservationService{
		ReservationRepo: reservationRepo,
		RestaurantRepo:  restaurantRepo,
		TableRepo:       tableRepo,
	}
}

// =====================================================
// CHECK RESTAURANT OWNER
// =====================================================

func (s *ReservationService) checkRestaurantOwner(
	restaurantID uuid.UUID,
	userID uuid.UUID,
) error {

	restaurant, err := s.RestaurantRepo.FindByID(
		restaurantID,
	)

	if err != nil {
		return errors.New("restaurant not found")
	}

	if restaurant.OwnerID != userID {
		return errors.New("you do not own this restaurant")
	}

	return nil
}

// =====================================================
// CREATE RESERVATION - CUSTOMER
// =====================================================

func (s *ReservationService) Create(
	userID uuid.UUID,
	reservation *models.Reservation,
) error {

	// Validate restaurant
	_, err := s.RestaurantRepo.FindByID(
		reservation.RestaurantID,
	)

	if err != nil {
		return errors.New("restaurant not found")
	}
	now := time.Now().UTC()
	reservationTime := reservation.ReservationTime.UTC()

	if !reservationTime.After(now) {
		return errors.New("reservation time must be in the future")
	}

	// Validate table
	table, err := s.TableRepo.FindByID(
		reservation.TableID,
	)

	if err != nil {
		return errors.New("table not found")
	}

	// Table must belong to selected restaurant
	if table.RestaurantID != reservation.RestaurantID {
		return errors.New(
			"table does not belong to this restaurant",
		)
	}

	// Check availability
	if !table.IsAvailable {
		return errors.New(
			"table is not available",
		)
	}

	// Validate guest count
	if reservation.GuestCount <= 0 {
		return errors.New(
			"guest count must be greater than zero",
		)
	}

	// Check capacity
	if reservation.GuestCount > table.Capacity {
		return errors.New(
			"guest count exceeds table capacity",
		)
	}

	// Check conflict
	conflict, err := s.ReservationRepo.HasConflict(
		reservation.TableID,
		reservation.ReservationTime,
	)
    if err != nil {
        return err
    }

    if conflict {
        return errors.New("table is already reserved for this time")
    }

	

	// Set authenticated user
	reservation.UserID = userID

	// Always set pending
	reservation.Status = "pending"

	return s.ReservationRepo.Create(reservation)
}

// =====================================================
// GET MY RESERVATIONS - CUSTOMER
// =====================================================

func (s *ReservationService) GetMyReservations(
	userID uuid.UUID,
) ([]models.Reservation, error) {

	return s.ReservationRepo.FindByUserID(
		userID,
	)
}

// =====================================================
// GET RESERVATION BY ID
// =====================================================

func (s *ReservationService) GetByID(
	id uuid.UUID,
) (*models.Reservation, error) {

	return s.ReservationRepo.FindByID(id)
}

// =====================================================
// CANCEL RESERVATION - CUSTOMER
// =====================================================

func (s *ReservationService) Cancel(
	reservationID uuid.UUID,
	userID uuid.UUID,
) (*models.Reservation, error) {

	reservation, err := s.ReservationRepo.FindByID(
		reservationID,
	)

	if err != nil {
		return nil, errors.New(
			"reservation not found",
		)
	}

	// Only reservation owner can cancel
	if reservation.UserID != userID {
		return nil, errors.New(
			"you do not own this reservation",
		)
	}

	if reservation.Status == "cancelled" {
		return nil, errors.New(
			"reservation already cancelled",
		)
	}

	if reservation.Status == "completed" {
		return nil, errors.New(
			"completed reservation cannot be cancelled",
		)
	}

	reservation.Status = "cancelled"

	if err := s.ReservationRepo.Update(
		reservation,
	); err != nil {
		return nil, err
	}

	return reservation, nil
}

// =====================================================
// GET RESTAURANT RESERVATIONS - OWNER
// =====================================================

func (s *ReservationService) GetRestaurantReservations(
	restaurantID uuid.UUID,
	userID uuid.UUID,
) ([]models.Reservation, error) {

	// Verify owner owns this restaurant
	if err := s.checkRestaurantOwner(
		restaurantID,
		userID,
	); err != nil {
		return nil, err
	}

	return s.ReservationRepo.FindByRestaurantID(
		restaurantID,
	)
}

// =====================================================
// CONFIRM RESERVATION - OWNER
// =====================================================

func (s *ReservationService) Confirm(
	reservationID uuid.UUID,
	userID uuid.UUID,
) (*models.Reservation, error) {

	reservation, err := s.ReservationRepo.FindByID(
		reservationID,
	)

	if err != nil {
		return nil, errors.New(
			"reservation not found",
		)
	}

	// Verify restaurant owner
	if err := s.checkRestaurantOwner(
		reservation.RestaurantID,
		userID,
	); err != nil {
		return nil, err
	}

	if reservation.Status != "pending" {
		return nil, errors.New(
			"only pending reservations can be confirmed",
		)
	}

	reservation.Status = "confirmed"

	if err := s.ReservationRepo.Update(
		reservation,
	); err != nil {
		return nil, err
	}

	return reservation, nil
}