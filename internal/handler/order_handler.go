package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/service"
)

type OrderHandler struct {
	Service *service.OrderService
}

func NewOrderHandler(
	service *service.OrderService,
) *OrderHandler {

	return &OrderHandler{
		Service: service,
	}
}
func (h *OrderHandler) Create(
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

	var order models.Order

	if err := json.NewDecoder(r.Body).
		Decode(&order); err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	if err := h.Service.Create(
		userID,
		&order,
	); err != nil {

		switch err.Error() {

		case "restaurant not found":
			writeError(
				w,
				http.StatusNotFound,
				err.Error(),
			)

		case "menu item not found":
			writeError(
				w,
				http.StatusNotFound,
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
		order,
	)
}
func (h *OrderHandler) GetMyOrders(
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

	orders, err := h.Service.GetMyOrders(
		userID,
	)

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"failed to get orders",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		orders,
	)
}

func (h *OrderHandler) GetByID(
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
			"invalid order id",
		)
		return
	}

	order, err := h.Service.GetByID(id)

	if err != nil {
		writeError(
			w,
			http.StatusNotFound,
			"order not found",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		order,
	)
}


func (h *OrderHandler) GetRestaurantOrders(
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

	orders, err := h.Service.GetRestaurantOrders(
		restaurantID,
		userID,
	)

	if err != nil {

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
				http.StatusInternalServerError,
				err.Error(),
			)
		}

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		orders,
	)
}

func (h *OrderHandler) UpdateStatus(
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

	orderID, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid order id",
		)
		return
	}

	var request struct {
		Status string `json:"status"`
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

	order, err := h.Service.UpdateStatus(
		orderID,
		userID,
		request.Status,
	)

	if err != nil {

		switch err.Error() {

		case "order not found",
			"restaurant not found":

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
		order,
	)
}

func (h *OrderHandler) Cancel(
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

	orderID, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid order id",
		)
		return
	}

	order, err := h.Service.Cancel(
		orderID,
		userID,
	)

	if err != nil {

		switch err.Error() {

		case "order not found":
			writeError(
				w,
				http.StatusNotFound,
				err.Error(),
			)

		case "you do not own this order":
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
		order,
	)
}