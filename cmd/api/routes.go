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

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			w.Write([]byte(
				`{"status":"available"}`,
			))
		},
	)

	// =========================
	// Authentication
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

	authMW := middleware.AuthMiddleware(
		app.AuthService,
	)

	mux.Handle(
		"/auth/me",
		authMW(
			http.HandlerFunc(
				app.AuthHandler.Me,
			),
		),
	)

	// =========================
	// Restaurants - Public
	// =========================

	mux.HandleFunc(
		"GET /api/restaurants",
		app.RestaurantHandler.GetAll,
	)

	mux.HandleFunc(
		"GET /api/restaurants/{id}",
		app.RestaurantHandler.GetByID,
	)

	// =========================
	// Restaurants - Protected
	// =========================

	mux.Handle(
		"/api/restaurants",
		authMW(
			http.HandlerFunc(
				app.RestaurantHandler.Create,
			),
		),
	)

	mux.Handle(
		"/api/my/restaurant",
		authMW(
			http.HandlerFunc(
				app.RestaurantHandler.MyRestaurant,
			),
		),
	)

	mux.Handle(
		"/api/restaurants/{id}/update",
		authMW(
			http.HandlerFunc(
				app.RestaurantHandler.Update,
			),
		),
	)

	mux.Handle(
		"/api/restaurants/{id}/delete",
		authMW(
			http.HandlerFunc(
				app.RestaurantHandler.Delete,
			),
		),
	)

	return mux
}