package main

import (
	"net/http"

	"restaurant-backend/internal/middleware"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	// =========================
	// Health
	// =========================

	mux.HandleFunc(
		"GET /healthcheck",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"available"}`))
		},
	)

	// =========================
	// Auth - Public
	// =========================

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

	// =========================
	// Auth Middleware
	// =========================

	authMW := middleware.AuthMiddleware(
		app.AuthService,
	)

	// =========================
	// Auth - Protected
	// =========================

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

	// Get all restaurants
	mux.HandleFunc(
		"GET /api/restaurants",
		app.RestaurantHandler.GetAll,
	)

	// Get restaurant by ID
	mux.HandleFunc(
		"GET /api/restaurants/{id}",
		app.RestaurantHandler.GetByID,
	)

	// ==================================================
	// RESTAURANT - PROTECTED
	// ==================================================

	// Create restaurant
	mux.Handle(
		"POST /api/restaurants",
		authMW(
			http.HandlerFunc(
				app.RestaurantHandler.Create,
			),
		),
	)

	// Get my restaurant
	mux.Handle(
		"GET /api/my/restaurant",
		authMW(
			middleware.RequireRole("owner", "admin")(
			http.HandlerFunc(
				app.RestaurantHandler.MyRestaurant,
			),),
		),
	)

	// Update restaurant
	mux.Handle(
		"PUT /api/restaurants/{id}",
		authMW(
			http.HandlerFunc(
				app.RestaurantHandler.Update,
			),
		),
	)

	// Delete restaurant
	mux.Handle(
		"DELETE /api/restaurants/{id}",
		authMW(
			http.HandlerFunc(
				app.RestaurantHandler.Delete,
			),
		),
	)

	// ==================================================
	// MENU - PUBLIC
	// ==================================================

	// Get all menu items for restaurant
	mux.HandleFunc(
		"GET /api/restaurants/{restaurantId}/menu",
		app.MenuHandler.GetRestaurantMenuPaginated,
	)

	// Get menu item by ID
	mux.HandleFunc(
		"GET /api/menu/{id}",
		app.MenuHandler.GetByID,
	)

	// ==================================================
	// MENU - PROTECTED
	// ==================================================

	// Create menu item
	mux.Handle(
		"POST /api/restaurants/{restaurantId}/menu",
		authMW(
			middleware.RequireRole("owner", "admin")(
			http.HandlerFunc(
				app.MenuHandler.Create,),
			),
		),
	)

	// Update menu item
	mux.Handle(
		"PUT /api/menu/{id}",
		authMW(
			middleware.RequireRole("owner", "admin")(
			http.HandlerFunc(
				app.MenuHandler.Update,
			),),
		),
	)

	// Update menu availability
	mux.Handle(
		"PATCH /api/menu/{id}/availability",
		authMW(
			middleware.RequireRole("owner", "admin")(
			http.HandlerFunc(
				app.MenuHandler.UpdateAvailability,
			),),
		),
	)

	// Delete menu item
	mux.Handle(
		"DELETE /api/menu/{id}",
		authMW(
			http.HandlerFunc(
				app.MenuHandler.Delete,
			),
		),
	)


	mux.HandleFunc(
	"GET /api/restaurants/{restaurantId}/tables",
	app.TableHandler.GetRestaurantTables,
)

mux.HandleFunc(
	"GET /api/tables/{id}",
	app.TableHandler.GetByID,
)

// =========================
// TABLE - PROTECTED
// =========================

mux.Handle(
	"POST /api/restaurants/{restaurantId}/tables",
	authMW(
		middleware.RequireRole("owner", "admin")(
			http.HandlerFunc(
				app.TableHandler.Create,
			),),
	),
)

mux.Handle(
	"PUT /api/tables/{id}",
	authMW(
		middleware.RequireRole("owner", "admin")(
			http.HandlerFunc(
				app.TableHandler.Update,
			),),
	),
)

mux.Handle(
	"PATCH /api/tables/{id}/availability",
	authMW(
		middleware.RequireRole("owner", "admin")(
			http.HandlerFunc(
				app.TableHandler.UpdateAvailability,
			),),
	),
)

mux.Handle(
	"DELETE /api/tables/{id}",
	authMW(
		middleware.RequireRole("owner", "admin")(
			http.HandlerFunc(
				app.TableHandler.Delete,
			),),
	),
)



// =========================
// RESERVATION - CUSTOMER
// =========================

mux.Handle(
	"POST /api/reservations",
	authMW(
		http.HandlerFunc(
			app.ReservationHandler.Create,
		),
	),
)

mux.Handle(
	"GET /api/my/reservations",
	authMW(
		http.HandlerFunc(
			app.ReservationHandler.GetMyReservations,
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
		http.HandlerFunc(
			app.ReservationHandler.Cancel,
		),
	),
)


// =========================
// RESERVATION - RESTAURANT OWNER
// =========================

mux.Handle(
	"GET /api/restaurants/{restaurantId}/reservations",
	authMW(
		http.HandlerFunc(
			app.ReservationHandler.GetRestaurantReservations,
		),
	),
)

mux.Handle(
	"PATCH /api/reservations/{id}/confirm",
	authMW(
		http.HandlerFunc(
			app.ReservationHandler.Confirm,
		),
	),
)


// =========================
// ORDER - CUSTOMER
// =========================

mux.Handle(
	"POST /api/orders",
	authMW(
	middleware.RequireRole("customer", "admin")(
			http.HandlerFunc(
				app.OrderHandler.Create,
			),),
	),
)

mux.Handle(
	"GET /api/my/orders",
	authMW(
		middleware.RequireRole("customer", "admin")(
			http.HandlerFunc(
				app.OrderHandler.GetMyOrders,
			),),
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



// =========================
// ORDER - RESTAURANT OWNER
// =========================

mux.Handle(
	"GET /api/restaurants/{restaurantId}/orders",
	authMW(
			middleware.RequireRole("owner", "admin")(
			http.HandlerFunc(
				app.OrderHandler.GetRestaurantOrders,
			),),
	),
)

mux.Handle(
	"PATCH /api/orders/{id}/status",
	authMW(
		middleware.RequireRole("owner", "admin")(
			http.HandlerFunc(
				app.OrderHandler.UpdateStatus,
			),),
	),
)



// =========================
// REVIEW
// =========================

// Get restaurant reviews - PUBLIC
mux.Handle(
	"GET /api/restaurants/{restaurantId}/reviews",
	http.HandlerFunc(
		app.ReviewHandler.GetRestaurantReviews,
	),
)

// Get review by ID - PUBLIC
mux.Handle(
	"GET /api/reviews/{id}",
	http.HandlerFunc(
		app.ReviewHandler.GetByID,
	),
)

// Create review - PROTECTED
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

// Get my reviews - PROTECTED
mux.Handle(
	"GET /api/my/reviews",
	authMW(
		http.HandlerFunc(
			app.ReviewHandler.GetMyReviews,
		),
	),
)

// Update review - PROTECTED
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

// Delete review - PROTECTED
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
	return mux
}