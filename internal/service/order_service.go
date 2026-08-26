package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"
)

type OrderService struct {
	OrderRepo      *repository.OrderRepository
	MenuRepo       *repository.MenuRepository
	RestaurantRepo *repository.RestaurantRepository
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	menuRepo *repository.MenuRepository,
	restaurantRepo *repository.RestaurantRepository,
) *OrderService {

	return &OrderService{
		OrderRepo:      orderRepo,
		MenuRepo:       menuRepo,
		RestaurantRepo: restaurantRepo,
	}
}



func (s *OrderService) Create(
	userID uuid.UUID,
	order *models.Order,
) error {

	// Check restaurant
	restaurant, err := s.RestaurantRepo.FindByID(
		order.RestaurantID,
	)

	if err != nil {
		return errors.New("restaurant not found")
	}

	_ = restaurant

	// Order must contain items
	if len(order.Items) == 0 {
		return errors.New("order must contain at least one item")
	}

	var total float64

	for i := range order.Items {

		item := &order.Items[i]

		if item.Quantity <= 0 {
			return errors.New(
				"quantity must be greater than zero",
			)
		}

		// Get menu item from DB
		menu, err := s.MenuRepo.FindByID(
			item.MenuItemID,
		)

		if err != nil {
			return errors.New(
				"menu item not found",
			)
		}

		// Make sure menu belongs to restaurant
		if menu.RestaurantID != order.RestaurantID {
			return errors.New(
				"menu item does not belong to this restaurant",
			)
		}

		// Check availability
		if !menu.IsAvailable {
			return errors.New(
				"menu item is not available: " + menu.Name,
			)
		}

		// IMPORTANT:
		// Never trust frontend price.
		item.Price = menu.Price

		total += item.Price * float64(item.Quantity)

		item.ID = uuid.Nil
	}

	order.UserID = userID
	order.TotalAmount = total
	order.Status = "pending"

	return s.OrderRepo.Create(order)
}
func (s *OrderService) GetByID(
	id uuid.UUID,
) (*models.Order, error) {

	return s.OrderRepo.FindByID(id)
}

func (s *OrderService) GetMyOrders(
	userID uuid.UUID,
) ([]models.Order, error) {

	return s.OrderRepo.FindByUserID(userID)
}
func (s *OrderService) checkRestaurantOwner(
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
		return errors.New(
			"you do not own this restaurant",
		)
	}

	return nil
}

func (s *OrderService) GetRestaurantOrders(
	restaurantID uuid.UUID,
	userID uuid.UUID,
) ([]models.Order, error) {

	if err := s.checkRestaurantOwner(
		restaurantID,
		userID,
	); err != nil {
		return nil, err
	}

	return s.OrderRepo.FindByRestaurantID(
		restaurantID,
	)
}

func (s *OrderService) UpdateStatus(
	orderID uuid.UUID,
	userID uuid.UUID,
	status string,
) (*models.Order, error) {

	order, err := s.OrderRepo.FindByID(
		orderID,
	)

	if err != nil {
		return nil, errors.New(
			"order not found",
		)
	}

	if err := s.checkRestaurantOwner(
		order.RestaurantID,
		userID,
	); err != nil {
		return nil, err
	}

	status = strings.ToLower(
		strings.TrimSpace(status),
	)

	validStatuses := map[string]bool{
		"pending":   true,
		"confirmed": true,
		"preparing": true,
		"ready":     true,
		"completed": true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return nil, errors.New(
			"invalid order status",
		)
	}

	order.Status = status

	if err := s.OrderRepo.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) Cancel(
	orderID uuid.UUID,
	userID uuid.UUID,
) (*models.Order, error) {

	order, err := s.OrderRepo.FindByID(
		orderID,
	)

	if err != nil {
		return nil, errors.New(
			"order not found",
		)
	}

	if order.UserID != userID {
		return nil, errors.New(
			"you do not own this order",
		)
	}

	if order.Status != "pending" {
		return nil, errors.New(
			"only pending orders can be cancelled",
		)
	}

	order.Status = "cancelled"

	if err := s.OrderRepo.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}