package service

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"
)

type RestaurantService struct {
	Repo *repository.RestaurantRepository
}

func NewRestaurantService(
	repo *repository.RestaurantRepository,
) *RestaurantService {
	return &RestaurantService{
		Repo: repo,
	}
}

func (s *RestaurantService) Create(
	ownerID uuid.UUID,
	restaurant *models.Restaurant,
) error {

	existing, err := s.Repo.FindByOwnerID(ownerID)

	if err == nil && existing != nil {
		return errors.New("owner already has a restaurant")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	restaurant.OwnerID = ownerID

	return s.Repo.Create(restaurant)
}

func (s *RestaurantService) GetByID(
	id uuid.UUID,
) (*models.Restaurant, error) {
	return s.Repo.FindByID(id)
}

func (s *RestaurantService) GetByOwner(
	ownerID uuid.UUID,
) (*models.Restaurant, error) {
	return s.Repo.FindByOwnerID(ownerID)
}

func (s *RestaurantService) GetAll() (
	[]models.Restaurant,
	error,
) {
	return s.Repo.FindAll()
}

func (s *RestaurantService) Update(
	ownerID uuid.UUID,
	id uuid.UUID,
	data *models.Restaurant,
) (*models.Restaurant, error) {

	restaurant, err := s.Repo.FindByID(id)

	if err != nil {
		return nil, err
	}

	if restaurant.OwnerID != ownerID {
		return nil, errors.New("you do not own this restaurant")
	}

	restaurant.Name = data.Name
	restaurant.Description = data.Description
	restaurant.Address = data.Address
	restaurant.City = data.City
	restaurant.State = data.State
	restaurant.Pincode = data.Pincode
	restaurant.Phone = data.Phone
	restaurant.Email = data.Email
	restaurant.Image = data.Image
	restaurant.OpeningTime = data.OpeningTime
	restaurant.ClosingTime = data.ClosingTime
	restaurant.IsOpen = data.IsOpen

	if err := s.Repo.Update(restaurant); err != nil {
		return nil, err
	}

	return restaurant, nil
}

func (s *RestaurantService) Delete(
	ownerID uuid.UUID,
	id uuid.UUID,
) error {

	restaurant, err := s.Repo.FindByID(id)

	if err != nil {
		return err
	}

	if restaurant.OwnerID != ownerID {
		return errors.New("you do not own this restaurant")
	}

	return s.Repo.Delete(id)
}