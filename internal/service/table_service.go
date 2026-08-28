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

	
	table.RestaurantID = restaurantID

	
	table.IsAvailable = true

	return s.TableRepo.Create(table)
}

func (s *TableService) GetRestaurantTables(
	restaurantID uuid.UUID,
) ([]models.RestaurantTable, error) {

	return s.TableRepo.FindByRestaurantID(
		restaurantID,
	)
}



func (s *TableService) GetByID(
	id uuid.UUID,
) (*models.RestaurantTable, error) {

	return s.TableRepo.FindByID(id)
}


func (s *TableService) Update(
	tableID uuid.UUID,
	userID uuid.UUID,
	data *models.RestaurantTable,
) (*models.RestaurantTable, error) {

	table, err := s.TableRepo.FindByID(tableID)

	if err != nil {
		return nil, errors.New("table not found")
	}

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

	if err := s.checkRestaurantOwner(
		table.RestaurantID,
		userID,
	); err != nil {
		return nil, err
	}


	table.IsAvailable = available


	if err := s.TableRepo.Update(table); err != nil {
		return nil, err
	}

	return table, nil
}