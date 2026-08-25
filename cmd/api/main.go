package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"restaurant-backend/internal/config"
	"restaurant-backend/internal/handler"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"
	"restaurant-backend/internal/service"
)

type Application struct {
	Config      *config.Config
	UserRepo    *repository.UserRepository
	RefreshRepo *repository.RefreshTokenRepository
	AuthService *service.AuthService
	AuthHandler *handler.AuthHandler
	RestaurantService *service.RestaurantService
	RestaurantHandler *handler.RestaurantHandler
}

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DBURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("database connected")

	if err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Restaurant{},
		&models.Category{},
		&models.MenuItem{},
		&models.RestaurantTable{},
		&models.Reservation{},
		&models.Order{},
		&models.OrderItem{},
		&models.Review{},
	); err != nil {
		log.Fatal(err)
	}

	log.Println("database migration completed")

	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	restaurantRepo := repository.NewRestaurantRepository(db)

restaurantService := service.NewRestaurantService(
	restaurantRepo,
)

restaurantHandler := handler.NewRestaurantHandler(
	restaurantService,
)

	authSvc := service.NewAuthService(
		userRepo,
		refreshTokenRepo,
		&cfg,
	)

	authHandler := handler.NewAuthHandler(authSvc)

	app := &Application{
		Config:      &cfg,
		UserRepo:    userRepo,
		RefreshRepo: refreshTokenRepo,
		AuthService: authSvc,
		AuthHandler: authHandler,
		RestaurantService: restaurantService,
	RestaurantHandler: restaurantHandler,
	}

	log.Println("server listening on :8080")

	if err := http.ListenAndServe(":8080", app.routes()); err != nil {
		log.Fatal(err)
	}
}