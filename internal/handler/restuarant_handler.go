package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/service"
)

type RestaurantHandler struct {
	Service *service.RestaurantService
}

func NewRestaurantHandler(
	service *service.RestaurantService,
) *RestaurantHandler {
	return &RestaurantHandler{
		Service: service,
	}
}

func (h *RestaurantHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := r.Context().
		Value(middleware.UserIDKey).
		(uuid.UUID)

	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var restaurant models.Restaurant

	if err := json.NewDecoder(r.Body).
		Decode(&restaurant); err != nil {

		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.Service.Create(
		userID,
		&restaurant,
	); err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		restaurant,
	)
}

func (h *RestaurantHandler) GetAll(
	w http.ResponseWriter,
	r *http.Request,
) {

	restaurants, err := h.Service.GetAll()

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"failed to get restaurants",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		restaurants,
	)
}

func (h *RestaurantHandler) GetByID(
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
			"invalid restaurant id",
		)
		return
	}

	restaurant, err := h.Service.GetByID(id)

	if err != nil {
		writeError(
			w,
			http.StatusNotFound,
			"restaurant not found",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		restaurant,
	)
}

func (h *RestaurantHandler) MyRestaurant(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := r.Context().
		Value(middleware.UserIDKey).
		(uuid.UUID)

	if !ok {
		writeError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
		)
		return
	}

	restaurant, err := h.Service.GetByOwner(userID)

	if err != nil {
		writeError(
			w,
			http.StatusNotFound,
			"restaurant not found",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		restaurant,
	)
}

func (h *RestaurantHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := r.Context().
		Value(middleware.UserIDKey).
		(uuid.UUID)

	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var data models.Restaurant

	if err := json.NewDecoder(r.Body).
		Decode(&data); err != nil {

		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	restaurant, err := h.Service.Update(
		userID,
		id,
		&data,
	)

	if err != nil {
		writeError(
			w,
			http.StatusForbidden,
			err.Error(),
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		restaurant,
	)
}

func (h *RestaurantHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := r.Context().
		Value(middleware.UserIDKey).
		(uuid.UUID)

	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.Service.Delete(
		userID,
		id,
	); err != nil {

		writeError(
			w,
			http.StatusForbidden,
			err.Error(),
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

