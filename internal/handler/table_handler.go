package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/service"
)

type TableHandler struct {
	Service *service.TableService
}

func NewTableHandler(
	service *service.TableService,
) *TableHandler {

	return &TableHandler{
		Service: service,
	}
}
func (h *TableHandler) Create(
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

	var table models.RestaurantTable

	if err := json.NewDecoder(r.Body).
		Decode(&table); err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	if err := h.Service.Create(
		restaurantID,
		userID,
		&table,
	); err != nil {

		switch err.Error() {

		case "restaurant not found":

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
		http.StatusCreated,
		table,
	)
}
func (h *TableHandler) GetRestaurantTables(
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

	tables, err := h.Service.GetRestaurantTables(
		restaurantID,
	)

	if err != nil {

		writeError(
			w,
			http.StatusInternalServerError,
			"failed to get tables",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		tables,
	)
}

func (h *TableHandler) GetByID(
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
			"invalid table id",
		)

		return
	}

	table, err := h.Service.GetByID(id)

	if err != nil {

		writeError(
			w,
			http.StatusNotFound,
			"table not found",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		table,
	)
}

func (h *TableHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	tableID, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid table id",
		)

		return
	}

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

	var data models.RestaurantTable

	if err := json.NewDecoder(r.Body).
		Decode(&data); err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)

		return
	}

	table, err := h.Service.Update(
		tableID,
		userID,
		&data,
	)

	if err != nil {

		switch err.Error() {

		case "table not found":

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
		table,
	)
}


func (h *TableHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	tableID, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid table id",
		)

		return
	}

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
		tableID,
		userID,
	)

	if err != nil {

		switch err.Error() {

		case "table not found":

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

func (h *TableHandler) UpdateAvailability(
	w http.ResponseWriter,
	r *http.Request,
) {

	tableID, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid table id",
		)

		return
	}

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

	var request struct {
		IsAvailable bool `json:"is_available"`
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

	table, err := h.Service.UpdateAvailability(
		tableID,
		userID,
		request.IsAvailable,
	)

	if err != nil {

		switch err.Error() {

		case "table not found":

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
		table,
	)
}