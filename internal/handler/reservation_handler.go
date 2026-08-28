package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/service"
)

type ReservationHandler struct {
	Service *service.ReservationService
}

func NewReservationHandler(
	service *service.ReservationService,
) *ReservationHandler {

	return &ReservationHandler{
		Service: service,
	}
}

// =====================================================
// CREATE RESERVATION
// POST /reservations
// =====================================================

func (h *ReservationHandler) Create(
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

	var reservation models.Reservation
	
    if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
        log.Printf("❌ JSON decode error: %v", err)
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    // 🔍 Debug: Check if time was parsed
    log.Printf("📥 Decoded reservation: %+v", reservation)
    log.Printf("📅 ReservationTime: %v (is zero? %v)",
		reservation.ReservationTime, reservation.ReservationTime.IsZero())

	if err := h.Service.Create(
		userID,
		&reservation,
	); err != nil {

		switch err.Error() {

		case "restaurant not found",
			"table not found":

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
		reservation,
	)
}

// =====================================================
// GET MY RESERVATIONS
// GET /reservations/my
// =====================================================

func (h *ReservationHandler) GetMyReservations(
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

	reservations, err := h.Service.GetMyReservations(
		userID,
	)

	if err != nil {

		writeError(
			w,
			http.StatusInternalServerError,
			"failed to get reservations",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		reservations,
	)
}

// =====================================================
// GET RESERVATION BY ID
// GET /reservations/{id}
// =====================================================

func (h *ReservationHandler) GetByID(
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
			"invalid reservation id",
		)
		return
	}

	reservation, err := h.Service.GetByID(id)

	if err != nil {

		writeError(
			w,
			http.StatusNotFound,
			"reservation not found",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		reservation,
	)
}

// =====================================================
// CANCEL RESERVATION
// PATCH /reservations/{id}/cancel
// =====================================================

func (h *ReservationHandler) Cancel(
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

	id, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid reservation id",
		)
		return
	}

	reservation, err := h.Service.Cancel(
		id,
		userID,
	)

	if err != nil {

		switch err.Error() {

		case "reservation not found":

			writeError(
				w,
				http.StatusNotFound,
				err.Error(),
			)

		case "you do not own this reservation":

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
		reservation,
	)
}

// =====================================================
// GET RESTAURANT RESERVATIONS - OWNER
// GET /restaurants/{restaurantId}/reservations
// =====================================================

func (h *ReservationHandler) GetRestaurantReservations(
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

	reservations, err := h.Service.GetRestaurantReservations(
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
		reservations,
	)
}

// =====================================================
// CONFIRM RESERVATION - OWNER
// PATCH /reservations/{id}/confirm
// =====================================================

func (h *ReservationHandler) Confirm(
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

	id, err := uuid.Parse(
		r.PathValue("id"),
	)

	if err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"invalid reservation id",
		)
		return
	}

	reservation, err := h.Service.Confirm(
		id,
		userID,
	)

	if err != nil {

		switch err.Error() {

		case "reservation not found",
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
		reservation,
	)
}