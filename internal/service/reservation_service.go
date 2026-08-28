package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"
)

var (
	ErrTableNotFound               = errors.New("table not found")
	ErrTableNotAvailable           = errors.New("table is not available")
	ErrGuestCountExceedsCapacity   = errors.New("guest count exceeds table capacity")
	ErrReservationTimePast         = errors.New("reservation time must be in the future")
	ErrTimeConflict                = errors.New("table is already reserved for an overlapping time")
	ErrNotReservationOwner         = errors.New("you do not own this reservation")
	ErrReservationAlreadyCancelled = errors.New("reservation already cancelled")
	ErrReservationCompleted        = errors.New("completed reservation cannot be cancelled")
	ErrInvalidStatusTransition     = errors.New("only pending reservations can be confirmed")
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
func (s *ReservationService) checkRestaurantOwner(restaurantID, userID uuid.UUID) error {
	restaurant, err := s.RestaurantRepo.FindByID(restaurantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRestaurantNotFound
		}
		return err
	}
	if restaurant.OwnerID != userID {
		return errors.New("you do not own this restaurant")
	}
	return nil
}



func (s *ReservationService) Create(userID uuid.UUID, reservation *models.Reservation) error {
	// 1. Validate restaurant
	_, err := s.RestaurantRepo.FindByID(reservation.RestaurantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRestaurantNotFound
		}
		return err
	}


	table, err := s.TableRepo.FindByID(reservation.TableID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTableNotFound
		}
		return err
	}

	if !table.IsAvailable {
		return ErrTableNotAvailable
	}

	if reservation.ReservationTime.Before(time.Now()) {
		return ErrReservationTimePast
	}

	
	if reservation.GuestCount <= 0 {
		return errors.New("guest count must be greater than zero")
	}

	if reservation.GuestCount > table.Capacity {
		return ErrGuestCountExceedsCapacity
	}

	
	conflict, err := s.ReservationRepo.HasConflict(
		reservation.TableID,
		reservation.ReservationTime,
	)
	if err != nil {
		return err
	}
	if conflict {
		return ErrTimeConflict
	}


	reservation.UserID = userID
	reservation.Status = "pending"

	return s.ReservationRepo.Create(reservation)
}



func (s *ReservationService) GetMyReservations(
	userID uuid.UUID,
) ([]models.Reservation, error) {

	return s.ReservationRepo.FindByUserID(
		userID,
	)
}



func (s *ReservationService) GetByID(
	id uuid.UUID,
) (*models.Reservation, error) {

	return s.ReservationRepo.FindByID(id)
}



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



func (s *ReservationService) GetRestaurantReservations(
	restaurantID uuid.UUID,
	userID uuid.UUID,
) ([]models.Reservation, error) {

	
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