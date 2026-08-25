package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Application struct {
	Logger *log.Logger
}

func main() {
	logger := log.Default()

	app := &Application{
		Logger: logger,
	}

	server := &http.Server{
		Addr:         ":8080",
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("server started on http://localhost:8080")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}
}

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", app.healthCheckHandler)
	mux.HandleFunc("GET /api/v1/restaurants", app.listRestaurantsHandler)

	return mux
}

func (app *Application) healthCheckHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	response := map[string]string{
		"status":  "ok",
		"service": "restaurant-api",
	}

	writeJSON(w, http.StatusOK, response)
}

func (app *Application) listRestaurantsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	response := map[string]any{
		"data":  []string{},
		"page":  1,
		"limit": 20,
	}

	writeJSON(w, http.StatusOK, response)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(jsonData); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
