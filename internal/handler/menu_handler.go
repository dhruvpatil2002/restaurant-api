package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/service"
)

type MenuHandler struct {
	Service *service.MenuService
}

func NewMenuHandler(
	service *service.MenuService,
) *MenuHandler {
	return &MenuHandler{
		Service: service,
	}
}

// =====================================================
// CREATE MENU
// POST /api/restaurants/{restaurantId}/menu
// =====================================================

func (h *MenuHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	// Get restaurant ID from URL
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

	// Get authenticated user ID from JWT middleware
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

	// Decode request body
	var menu models.Menu

	if err := json.NewDecoder(r.Body).Decode(&menu); err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)

		return
	}

	// Create menu
	if err := h.Service.Create(
		restaurantID,
		userID,
		&menu,
	); err != nil {

		if err.Error() == "you do not own this restaurant" {
			writeError(
				w,
				http.StatusForbidden,
				err.Error(),
			)
			return
		}

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
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		menu,
	)
}

// =====================================================
// GET RESTAURANT MENU
// GET /api/restaurants/{restaurantId}/menu
// =====================================================

func (h *MenuHandler) GetRestaurantMenu(
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

	menus, err := h.Service.GetRestaurantMenu(
		restaurantID,
	)

	if err != nil {

		writeError(
			w,
			http.StatusInternalServerError,
			"failed to get menu",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		menus,
	)
}

// =====================================================
// GET MENU BY ID
// GET /api/menu/{id}
// =====================================================

func (h *MenuHandler) GetByID(
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
			"invalid menu id",
		)

		return
	}

	menu, err := h.Service.GetByID(id)

	if err != nil {

		writeError(
			w,
			http.StatusNotFound,
			"menu item not found",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		menu,
	)
}

// =====================================================
// UPDATE MENU
// PUT /api/menu/{id}
// =====================================================

func (h *MenuHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	menuID, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid menu id",
		)

		return
	}

	// Get authenticated user
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

	var data models.Menu

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)

		return
	}

	menu, err := h.Service.Update(
		menuID,
		userID,
		&data,
	)

	if err != nil {

		switch err.Error() {

		case "menu item not found":
			writeError(
				w,
				http.StatusNotFound,
				err.Error(),
			)

		case "you do not own this restaurant":
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
		menu,
	)
}

// =====================================================
// DELETE MENU
// DELETE /api/menu/{id}
// =====================================================

func (h *MenuHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	menuID, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid menu id",
		)

		return
	}

	// Get authenticated user
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

	err = h.Service.Delete(
		menuID,
		userID,
	)

	if err != nil {

		switch err.Error() {

		case "menu item not found":
			writeError(
				w,
				http.StatusNotFound,
				err.Error(),
			)

		case "you do not own this restaurant":
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

// =====================================================
// UPDATE AVAILABILITY
// PATCH /api/menu/{id}/availability
// =====================================================

func (h *MenuHandler) UpdateAvailability(
	w http.ResponseWriter,
	r *http.Request,
) {

	menuID, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid menu id",
		)

		return
	}

	// Get authenticated user
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

	// Request body
	var request struct {
		IsAvailable bool `json:"is_available"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)

		return
	}

	menu, err := h.Service.UpdateAvailability(
		menuID,
		userID,
		request.IsAvailable,
	)

	if err != nil {

		switch err.Error() {

		case "menu item not found":
			writeError(
				w,
				http.StatusNotFound,
				err.Error(),
			)

		case "you do not own this restaurant":
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
		menu,
	)
}