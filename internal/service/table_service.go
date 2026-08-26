package service

import (
	"errors"

	"github.com/google/uuid"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"
)

type TableService struct {
	TableRepo      *repository.TableRepository
	RestaurantRepo *repository.RestaurantRepository
}

func NewTableService(
	tableRepo *repository.TableRepository,
	restaurantRepo *repository.RestaurantRepository,
) *TableService {

	return &TableService{
		TableRepo:      tableRepo,
		RestaurantRepo: restaurantRepo,
	}
}

// =====================================================
// CHECK RESTAURANT OWNER
// =====================================================

func (s *TableService) checkRestaurantOwner(
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
// CREATE TABLE
// =====================================================

func (s *TableService) Create(
	restaurantID uuid.UUID,
	userID uuid.UUID,
	table *models.RestaurantTable,
) error {

	// Check restaurant exists and user owns it
	if err := s.checkRestaurantOwner(
		restaurantID,
		userID,
	); err != nil {
		return err
	}

	if table.TableNumber <= 0 {
		return errors.New(
			"table number must be greater than zero",
		)
	}

	if table.Capacity <= 0 {
		return errors.New(
			"capacity must be greater than zero",
		)
	}

	// Never trust restaurant_id from request body
	table.RestaurantID = restaurantID

	// New tables are available by default
	table.IsAvailable = true

	return s.TableRepo.Create(table)
}

// =====================================================
// GET RESTAURANT TABLES
// =====================================================

func (s *TableService) GetRestaurantTables(
	restaurantID uuid.UUID,
) ([]models.RestaurantTable, error) {

	return s.TableRepo.FindByRestaurantID(
		restaurantID,
	)
}

// =====================================================
// GET TABLE BY ID
// =====================================================

func (s *TableService) GetByID(
	id uuid.UUID,
) (*models.RestaurantTable, error) {

	return s.TableRepo.FindByID(id)
}

// =====================================================
// UPDATE TABLE
// =====================================================

func (s *TableService) Update(
	tableID uuid.UUID,
	userID uuid.UUID,
	data *models.RestaurantTable,
) (*models.RestaurantTable, error) {

	// Find existing table
	table, err := s.TableRepo.FindByID(tableID)

	if err != nil {
		return nil, errors.New("table not found")
	}

	// Check ownership
	if err := s.checkRestaurantOwner(
		table.RestaurantID,
		userID,
	); err != nil {
		return nil, err
	}

	if data.TableNumber <= 0 {
		return nil, errors.New(
			"table number must be greater than zero",
		)
	}

	if data.Capacity <= 0 {
		return nil, errors.New(
			"capacity must be greater than zero",
		)
	}

	// Update only editable fields
	table.TableNumber = data.TableNumber
	table.Capacity = data.Capacity

	if err := s.TableRepo.Update(table); err != nil {
		return nil, err
	}

	return table, nil
}

// =====================================================
// DELETE TABLE
// =====================================================

func (s *TableService) Delete(
	tableID uuid.UUID,
	userID uuid.UUID,
) error {

	table, err := s.TableRepo.FindByID(tableID)

	if err != nil {
		return errors.New("table not found")
	}

	// Check ownership
	if err := s.checkRestaurantOwner(
		table.RestaurantID,
		userID,
	); err != nil {
		return err
	}

	return s.TableRepo.Delete(tableID)
}

// =====================================================
// UPDATE AVAILABILITY
// =====================================================

func (s *TableService) UpdateAvailability(
	tableID uuid.UUID,
	userID uuid.UUID,
	available bool,
) (*models.RestaurantTable, error) {

	// Find table
	table, err := s.TableRepo.FindByID(tableID)

	if err != nil {
		return nil, errors.New("table not found")
	}

	// Check restaurant ownership
	if err := s.checkRestaurantOwner(
		table.RestaurantID,
		userID,
	); err != nil {
		return nil, err
	}

	// Update availability
	table.IsAvailable = available

	// Save
	if err := s.TableRepo.Update(table); err != nil {
		return nil, err
	}

	return table, nil
}