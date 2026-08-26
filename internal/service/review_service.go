package service

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"
)

type ReviewService struct {
	ReviewRepo     *repository.ReviewRepository
	RestaurantRepo *repository.RestaurantRepository
	OrderRepo      *repository.OrderRepository
}

func NewReviewService(
	reviewRepo *repository.ReviewRepository,
	restaurantRepo *repository.RestaurantRepository,
	orderRepo *repository.OrderRepository,
) *ReviewService {

	return &ReviewService{
		ReviewRepo:     reviewRepo,
		RestaurantRepo: restaurantRepo,
		OrderRepo:      orderRepo,
	}
}

func (s *ReviewService) hasCompletedOrder(
	userID uuid.UUID,
	restaurantID uuid.UUID,
) bool {

	var orders []models.Order

	orders, err := s.OrderRepo.FindByUserID(userID)

	if err != nil {
		return false
	}

	for _, order := range orders {

		if order.RestaurantID == restaurantID &&
			order.Status == "completed" {

			return true
		}
	}

	return false
}

func (s *ReviewService) Create(
	userID uuid.UUID,
	review *models.Review,
) error {

	// Check restaurant exists
	_, err := s.RestaurantRepo.FindByID(
		review.RestaurantID,
	)

	if err != nil {
		return errors.New(
			"restaurant not found",
		)
	}

	// Validate rating
	if review.Rating < 1 || review.Rating > 5 {
		return errors.New(
			"rating must be between 1 and 5",
		)
	}

	// Check completed order
	if !s.hasCompletedOrder(
		userID,
		review.RestaurantID,
	) {
		return errors.New(
			"you can review only restaurants you have completed an order from",
		)
	}

	// Check duplicate review
	existing, err := s.ReviewRepo.FindByUserAndRestaurant(
		userID,
		review.RestaurantID,
	)

	if err == nil && existing != nil {
		return errors.New(
			"you have already reviewed this restaurant",
		)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	review.UserID = userID

	return s.ReviewRepo.Create(review)
}
func (s *ReviewService) GetRestaurantReviews(
	restaurantID uuid.UUID,
) ([]models.Review, error) {

	_, err := s.RestaurantRepo.FindByID(
		restaurantID,
	)

	if err != nil {
		return nil, errors.New(
			"restaurant not found",
		)
	}

	return s.ReviewRepo.FindByRestaurantID(
		restaurantID,
	)
}
func (s *ReviewService) GetMyReviews(
	userID uuid.UUID,
) ([]models.Review, error) {

	return s.ReviewRepo.FindByUserID(
		userID,
	)
}
func (s *ReviewService) GetByID(
	id uuid.UUID,
) (*models.Review, error) {

	return s.ReviewRepo.FindByID(id)
}

func (s *ReviewService) Update(
	reviewID uuid.UUID,
	userID uuid.UUID,
	data *models.Review,
) (*models.Review, error) {

	review, err := s.ReviewRepo.FindByID(
		reviewID,
	)

	if err != nil {
		return nil, errors.New(
			"review not found",
		)
	}

	if review.UserID != userID {
		return nil, errors.New(
			"you do not own this review",
		)
	}

	if data.Rating < 1 || data.Rating > 5 {
		return nil, errors.New(
			"rating must be between 1 and 5",
		)
	}

	review.Rating = data.Rating
	review.Comment = data.Comment

	if err := s.ReviewRepo.Update(review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *ReviewService) Delete(
	reviewID uuid.UUID,
	userID uuid.UUID,
) error {

	review, err := s.ReviewRepo.FindByID(
		reviewID,
	)

	if err != nil {
		return errors.New(
			"review not found",
		)
	}

	if review.UserID != userID {
		return errors.New(
			"you do not own this review",
		)
	}

	return s.ReviewRepo.Delete(
		reviewID,
	)
}