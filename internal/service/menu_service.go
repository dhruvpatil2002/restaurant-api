package service

import (
	"errors"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"
)

type MenuService struct {
	MenuRepo       *repository.MenuRepository
	RestaurantRepo *repository.RestaurantRepository
}

func (s *MenuService) GetRestaurantMenu(restaurantID uuid.UUID) ([]models.Menu, error) {
	return s.MenuRepo.FindByRestaurantID(restaurantID)
}

func NewMenuService(
	menuRepo *repository.MenuRepository,
	restaurantRepo *repository.RestaurantRepository,
) *MenuService {
	return &MenuService{
		MenuRepo:       menuRepo,
		RestaurantRepo: restaurantRepo,
	}
}

// Create
func (s *MenuService) Create(
	restaurantID uuid.UUID,
	userID uuid.UUID,
	menu *models.Menu,
) error {

	if err := s.checkRestaurantOwner(restaurantID, userID); err != nil {
		return err
	}

	menu.RestaurantID = restaurantID

	if menu.Name == "" {
		return errors.New("menu name is required")
	}

	if menu.Category == "" {
		return errors.New("category is required")
	}

	if menu.Price <= 0 {
		return errors.New("price must be greater than zero")
	}

	return s.MenuRepo.Create(menu)
}

// Get by ID
func (s *MenuService) GetByID(
	id uuid.UUID,
) (*models.Menu, error) {

	return s.MenuRepo.FindByID(id)
}



// Get by category
func (s *MenuService) GetByCategory(
	restaurantID uuid.UUID,
	category string,
) ([]models.Menu, error) {

	return s.MenuRepo.FindByCategory(
		restaurantID,
		category,
	)
}

func (s *MenuService) checkRestaurantOwner(
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

// Update
func (s *MenuService) Update(
	menuID uuid.UUID,
	userID uuid.UUID,
	data *models.Menu,
) (*models.Menu, error) {

	menu, err := s.MenuRepo.FindByID(menuID)

	if err != nil {
		return nil, errors.New("menu item not found")
	}

	if err := s.checkRestaurantOwner(
		menu.RestaurantID,
		userID,
	); err != nil {
		return nil, err
	}

	if data.Name == "" {
		return nil, errors.New("menu name is required")
	}

	if data.Category == "" {
		return nil, errors.New("category is required")
	}

	if data.Price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}

	menu.Name = data.Name
	menu.Category = data.Category
	menu.Description = data.Description
	menu.Price = data.Price
	menu.Image = data.Image

	if err := s.MenuRepo.Update(menu); err != nil {
		return nil, err
	}

	return menu, nil
}

// Delete
func (s *MenuService) Delete(
	menuID uuid.UUID,
	userID uuid.UUID,
) error {

	menu, err := s.MenuRepo.FindByID(menuID)

	if err != nil {
		return errors.New("menu item not found")
	}

	if err := s.checkRestaurantOwner(
		menu.RestaurantID,
		userID,
	); err != nil {
		return err
	}

	return s.MenuRepo.Delete(menuID)
}

// Availability
func (s *MenuService) UpdateAvailability(
	menuID uuid.UUID,
	userID uuid.UUID,
	available bool,
) (*models.Menu, error) {

	menu, err := s.MenuRepo.FindByID(menuID)

	if err != nil {
		return nil, errors.New("menu item not found")
	}

	if err := s.checkRestaurantOwner(
		menu.RestaurantID,
		userID,
	); err != nil {
		return nil, err
	}

	menu.IsAvailable = available

	if err := s.MenuRepo.Update(menu); err != nil {
		return nil, err
	}

	return menu, nil
}

// Check menu exists
func (s *MenuService) Exists(
	id uuid.UUID,
) bool {

	_, err := s.MenuRepo.FindByID(id)

	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (s *MenuService) GetRestaurantMenuPaginated(
	restaurantID uuid.UUID,
	page int,
	limit int,
	search string,
	category string,
	available *bool,
) (*models.PaginatedMenus, error) {

	// Default page
	if page < 1 {
		page = 1
	}

	// Default limit
	if limit < 1 {
		limit = 10
	}

	// Maximum limit
	if limit > 100 {
		limit = 100
	}

	menus, total, err := s.MenuRepo.FindPaginated(
		restaurantID,
		page,
		limit,
		search,
		category,
		available,
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

	return &models.PaginatedMenus{
		Data: menus,
		Pagination: models.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}