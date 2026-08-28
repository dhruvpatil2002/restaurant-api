package main

import (
	"net/http"

	"restaurant-backend/internal/middleware"
)

func (app *Application) routes() http.Handler {

	mux := http.NewServeMux()

	// ==================================================
	// HEALTH
	// ==================================================

	mux.HandleFunc(
		"GET /healthcheck",
		func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusOK)

			_, _ = w.Write(
				[]byte(`{"status":"available"}`),
			)
		},
	)

	// ==================================================
	// AUTH - PUBLIC
	// ==================================================

	mux.HandleFunc(
		"POST /auth/register",
		app.AuthHandler.Register,
	)

	mux.HandleFunc(
		"POST /auth/login",
		app.AuthHandler.Login,
	)

	mux.HandleFunc(
		"POST /auth/refresh",
		app.AuthHandler.Refresh,
	)

	mux.HandleFunc(
		"POST /auth/logout",
		app.AuthHandler.Logout,
	)

	// ==================================================
	// AUTH MIDDLEWARE
	// ==================================================

	authMW := middleware.AuthMiddleware(
		app.AuthService,
	)

	// ==================================================
	// AUTH - PROTECTED
	// ==================================================

	mux.Handle(
		"GET /auth/me",
		authMW(
			http.HandlerFunc(
				app.AuthHandler.Me,
			),
		),
	)

	// ==================================================
	// RESTAURANT - PUBLIC
	// ==================================================

	mux.HandleFunc(
		"GET /api/restaurants",
		app.RestaurantHandler.GetAll,
	)

	mux.HandleFunc(
		"GET /api/restaurants/{id}",
		app.RestaurantHandler.GetByID,
	)

	// ==================================================
	// RESTAURANT - OWNER
	// ==================================================

	mux.Handle(
		"POST /api/restaurants",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.RestaurantHandler.Create,
				),
			),
		),
	)

	mux.Handle(
		"GET /api/my/restaurant",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.RestaurantHandler.MyRestaurant,
				),
			),
		),
	)

	mux.Handle(
		"PUT /api/restaurants/{id}",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.RestaurantHandler.Update,
				),
			),
		),
	)

	mux.Handle(
		"DELETE /api/restaurants/{id}",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.RestaurantHandler.Delete,
				),
			),
		),
	)

	// ==================================================
	// MENU - PUBLIC
	// ==================================================

	mux.HandleFunc(
		"GET /api/restaurants/{restaurantId}/menu",
		app.MenuHandler.GetRestaurantMenuPaginated,
	)

	mux.HandleFunc(
		"GET /api/menu/{id}",
		app.MenuHandler.GetByID,
	)

	// ==================================================
	// MENU - OWNER
	// ==================================================

	mux.Handle(
		"POST /api/restaurants/{restaurantId}/menu",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.MenuHandler.Create,
				),
			),
		),
	)

	mux.Handle(
		"PUT /api/menu/{id}",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.MenuHandler.Update,
				),
			),
		),
	)

	mux.Handle(
		"PATCH /api/menu/{id}/availability",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.MenuHandler.UpdateAvailability,
				),
			),
		),
	)

	mux.Handle(
		"DELETE /api/menu/{id}",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.MenuHandler.Delete,
				),
			),
		),
	)

	// ==================================================
	// TABLE - PUBLIC
	// ==================================================

	mux.HandleFunc(
		"GET /api/restaurants/{restaurantId}/tables",
		app.TableHandler.GetRestaurantTables,
	)

	mux.HandleFunc(
		"GET /api/tables/{id}",
		app.TableHandler.GetByID,
	)

	// ==================================================
	// TABLE - OWNER
	// ==================================================

	mux.Handle(
		"POST /api/restaurants/{restaurantId}/tables",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.TableHandler.Create,
				),
			),
		),
	)

	mux.Handle(
		"PUT /api/tables/{id}",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.TableHandler.Update,
				),
			),
		),
	)

	mux.Handle(
		"PATCH /api/tables/{id}/availability",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.TableHandler.UpdateAvailability,
				),
			),
		),
	)

	mux.Handle(
		"DELETE /api/tables/{id}",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.TableHandler.Delete,
				),
			),
		),
	)

	// ==================================================
	// RESERVATION - CUSTOMER
	// ==================================================

	mux.Handle(
		"POST /api/reservations",
		authMW(
			middleware.RequireRole("customer", "admin")(
				http.HandlerFunc(
					app.ReservationHandler.Create,
				),
			),
		),
	)

	mux.Handle(
		"GET /api/my/reservations",
		authMW(
			middleware.RequireRole("customer", "admin")(
				http.HandlerFunc(
					app.ReservationHandler.GetMyReservations,
				),
			),
		),
	)

	mux.Handle(
		"GET /api/reservations/{id}",
		authMW(
			http.HandlerFunc(
				app.ReservationHandler.GetByID,
			),
		),
	)

	mux.Handle(
		"PATCH /api/reservations/{id}/cancel",
		authMW(
			middleware.RequireRole("customer", "admin")(
				http.HandlerFunc(
					app.ReservationHandler.Cancel,
				),
			),
		),
	)

	// ==================================================
	// RESERVATION - OWNER
	// ==================================================

	mux.Handle(
		"GET /api/restaurants/{restaurantId}/reservations",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.ReservationHandler.GetRestaurantReservations,
				),
			),
		),
	)

	mux.Handle(
		"PATCH /api/reservations/{id}/confirm",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.ReservationHandler.Confirm,
				),
			),
		),
	)

	// ==================================================
	// ORDER - CUSTOMER
	// ==================================================

	mux.Handle(
		"POST /api/orders",
		authMW(
			middleware.RequireRole("customer", "admin")(
				http.HandlerFunc(
					app.OrderHandler.Create,
				),
			),
		),
	)

	mux.Handle(
		"GET /api/my/orders",
		authMW(
			middleware.RequireRole("customer", "admin")(
				http.HandlerFunc(
					app.OrderHandler.GetMyOrders,
				),
			),
		),
	)

	mux.Handle(
		"GET /api/orders/{id}",
		authMW(
			http.HandlerFunc(
				app.OrderHandler.GetByID,
			),
		),
	)

	mux.Handle(
		"PATCH /api/orders/{id}/cancel",
		authMW(
			middleware.RequireRole("customer", "admin")(
				http.HandlerFunc(
					app.OrderHandler.Cancel,
				),
			),
		),
	)

	// ==================================================
	// ORDER - OWNER
	// ==================================================

	mux.Handle(
		"GET /api/restaurants/{restaurantId}/orders",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.OrderHandler.GetRestaurantOrders,
				),
			),
		),
	)

	mux.Handle(
		"PATCH /api/orders/{id}/status",
		authMW(
			middleware.RequireRole("owner", "admin")(
				http.HandlerFunc(
					app.OrderHandler.UpdateStatus,
				),
			),
		),
	)

	// ==================================================
	// REVIEW - PUBLIC
	// ==================================================

	mux.Handle(
		"GET /api/restaurants/{restaurantId}/reviews",
		http.HandlerFunc(
			app.ReviewHandler.GetRestaurantReviews,
		),
	)

	mux.Handle(
		"GET /api/reviews/{id}",
		http.HandlerFunc(
			app.ReviewHandler.GetByID,
		),
	)

	// ==================================================
	// REVIEW - CUSTOMER
	// ==================================================

	mux.Handle(
		"POST /api/restaurants/{restaurantId}/reviews",
		authMW(
			middleware.RequireRole("customer", "admin")(
				http.HandlerFunc(
					app.ReviewHandler.Create,
				),
			),
		),
	)

	mux.Handle(
		"GET /api/my/reviews",
		authMW(
			middleware.RequireRole("customer", "admin")(
				http.HandlerFunc(
					app.ReviewHandler.GetMyReviews,
				),
			),
		),
	)

	mux.Handle(
		"PUT /api/reviews/{id}",
		authMW(
			middleware.RequireRole("customer", "admin")(
				http.HandlerFunc(
					app.ReviewHandler.Update,
				),
			),
		),
	)

	mux.Handle(
		"DELETE /api/reviews/{id}",
		authMW(
			middleware.RequireRole("customer", "admin")(
				http.HandlerFunc(
					app.ReviewHandler.Delete,
				),
			),
		),
	)

	// ==================================================
	// CORS
	// ==================================================

	corsHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set(
				"Access-Control-Allow-Origin",
				"*",
			)

			w.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, PATCH, DELETE, OPTIONS",
			)

			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization",
			)

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			mux.ServeHTTP(w, r)
		},
	)

	// ==================================================
	// LOGGING
	// ==================================================

	return middleware.Logging(
		corsHandler,
	)
}