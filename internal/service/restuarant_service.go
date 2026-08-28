package service

import (
	"errors"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"
)

var (
	ErrRestaurantNotFound      = errors.New("restaurant not found")
	ErrOwnerAlreadyHasRest     = errors.New("owner already has a restaurant")
	ErrNotRestaurantOwner      = errors.New("you do not own this restaurant")
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

// =====================================================
// CREATE
// =====================================================

func (s *RestaurantService) Create(
	ownerID uuid.UUID,
	restaurant *models.Restaurant,
) error {

	existing, err := s.Repo.FindByOwnerID(ownerID)

	if err == nil && existing != nil {
		return ErrOwnerAlreadyHasRest
	}

	if err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {

		return err
	}

	// NEVER accept OwnerID from client
	restaurant.OwnerID = ownerID

	return s.Repo.Create(restaurant)
}

// =====================================================
// GET BY OWNER
// =====================================================

func (s *RestaurantService) GetByOwner(
	ownerID uuid.UUID,
) (*models.Restaurant, error) {

	restaurant, err := s.Repo.FindByOwnerID(ownerID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRestaurantNotFound
	}

	return restaurant, err
}

// =====================================================
// GET ALL
// =====================================================

func (s *RestaurantService) GetAll() (
	[]models.Restaurant,
	error,
) {

	return s.Repo.FindAll()
}

// =====================================================
// GET ALL PAGINATED
// =====================================================

func (s *RestaurantService) GetAllPaginated(
	page int,
	limit int,
	search string,
) (*models.PaginatedRestaurants, error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	restaurants, total, err :=
		s.Repo.FindAllPaginated(
			page,
			limit,
			search,
		)

	if err != nil {
		return nil, err
	}

	totalPages := 0

	if total > 0 {
		totalPages = int(
			math.Ceil(
				float64(total)/float64(limit),
			),
		)
	}

	return &models.PaginatedRestaurants{
		Data: restaurants,
		Pagination: models.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// =====================================================
// GET BY ID
// =====================================================

func (s *RestaurantService) GetByID(
	id uuid.UUID,
) (*models.Restaurant, error) {

	restaurant, err := s.Repo.FindByID(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRestaurantNotFound
	}

	return restaurant, err
}

// =====================================================
// UPDATE
// =====================================================

func (s *RestaurantService) Update(
	ownerID uuid.UUID,
	id uuid.UUID,
	data *models.Restaurant,
) (*models.Restaurant, error) {

	restaurant, err := s.Repo.FindByID(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRestaurantNotFound
	}

	if err != nil {
		return nil, err
	}

	if restaurant.OwnerID != ownerID {
		return nil, ErrNotRestaurantOwner
	}

	// Update allowed fields only
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

// =====================================================
// DELETE
// =====================================================

func (s *RestaurantService) Delete(
	ownerID uuid.UUID,
	id uuid.UUID,
) error {

	restaurant, err := s.Repo.FindByID(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrRestaurantNotFound
	}

	if err != nil {
		return err
	}

	if restaurant.OwnerID != ownerID {
		return ErrNotRestaurantOwner
	}

	return s.Repo.Delete(id)
}