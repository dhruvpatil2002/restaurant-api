package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/service"
)

type ReviewHandler struct {
	Service *service.ReviewService
}

func NewReviewHandler(
	service *service.ReviewService,
) *ReviewHandler {

	return &ReviewHandler{
		Service: service,
	}
}

func (h *ReviewHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	userIDValue := r.Context().Value(
		middleware.UserIDKey,
	)

	userID, ok := userIDValue.(uuid.UUID)

	if !ok {
		writeError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
		)
		return
	}

	restaurantID, err := uuid.Parse(
		r.PathValue("restaurantId"),
	)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid restaurant id",
		)
		return
	}

	var request struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).
		Decode(&request); err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	review := models.Review{
		RestaurantID: restaurantID,
		Rating:       request.Rating,
		Comment:      request.Comment,
	}

	if err := h.Service.Create(
		userID,
		&review,
	); err != nil {

		switch err.Error() {

		case "restaurant not found":
			writeError(
				w,
				http.StatusNotFound,
				err.Error(),
			)

		case "you have already reviewed this restaurant":
			writeError(
				w,
				http.StatusConflict,
				err.Error(),
			)

		case "you can review only restaurants you have completed an order from":
			writeError(
				w,
				http.StatusForbidden,
				err.Error(),
			)

		default:
			writeError(
				w,
				http.StatusBadRequest,
				err.Error(),
			)
		}

		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		review,
	)
}


func (h *ReviewHandler) GetRestaurantReviews(
	w http.ResponseWriter,
	r *http.Request,
) {

	restaurantID, err := uuid.Parse(
		r.PathValue("restaurantId"),
	)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid restaurant id",
		)
		return
	}

	reviews, err := h.Service.GetRestaurantReviews(
		restaurantID,
	)

	if err != nil {

		if err.Error() == "restaurant not found" {
			writeError(
				w,
				http.StatusNotFound,
				err.Error(),
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"failed to get reviews",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		reviews,
	)
}

func (h *ReviewHandler) GetMyReviews(
	w http.ResponseWriter,
	r *http.Request,
) {

	userIDValue := r.Context().Value(
		middleware.UserIDKey,
	)

	userID, ok := userIDValue.(uuid.UUID)

	if !ok {
		writeError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
		)
		return
	}

	reviews, err := h.Service.GetMyReviews(
		userID,
	)

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"failed to get reviews",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		reviews,
	)
}

func (h *ReviewHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid review id",
		)
		return
	}

	review, err := h.Service.GetByID(id)

	if err != nil {
		writeError(
			w,
			http.StatusNotFound,
			"review not found",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		review,
	)
}

func (h *ReviewHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	userIDValue := r.Context().Value(
		middleware.UserIDKey,
	)

	userID, ok := userIDValue.(uuid.UUID)

	if !ok {
		writeError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
		)
		return
	}

	reviewID, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid review id",
		)
		return
	}

	var request struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).
		Decode(&request); err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	review, err := h.Service.Update(
		reviewID,
		userID,
		&models.Review{
			Rating:  request.Rating,
			Comment: request.Comment,
		},
	)

	if err != nil {

		switch err.Error() {

		case "review not found":
			writeError(
				w,
				http.StatusNotFound,
				err.Error(),
			)

		case "you do not own this review":
			writeError(
				w,
				http.StatusForbidden,
				err.Error(),
			)

		default:
			writeError(
				w,
				http.StatusBadRequest,
				err.Error(),
			)
		}

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		review,
	)
}

func (h *ReviewHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	userIDValue := r.Context().Value(
		middleware.UserIDKey,
	)

	userID, ok := userIDValue.(uuid.UUID)

	if !ok {
		writeError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
		)
		return
	}

	reviewID, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid review id",
		)
		return
	}

	if err := h.Service.Delete(
		reviewID,
		userID,
	); err != nil {

		switch err.Error() {

		case "review not found":
			writeError(
				w,
				http.StatusNotFound,
				err.Error(),
			)

		case "you do not own this review":
			writeError(
				w,
				http.StatusForbidden,
				err.Error(),
			)

		default:
			writeError(
				w,
				http.StatusBadRequest,
				err.Error(),
			)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}