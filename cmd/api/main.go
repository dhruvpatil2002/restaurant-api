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

	MenuService *service.MenuService
	MenuHandler *handler.MenuHandler
	TableService *service.TableService
TableHandler *handler.TableHandler
ReservationService *service.ReservationService
ReservationHandler *handler.ReservationHandler
OrderService *service.OrderService
OrderHandler *handler.OrderHandler
ReviewService *service.ReviewService
ReviewHandler *handler.ReviewHandler
}

func main() {
	_ = godotenv.Load()

	// =========================
	// CONFIG
	// =========================

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// =========================
	// DATABASE
	// =========================

	db, err := gorm.Open(
		postgres.Open(cfg.DBURL),
		&gorm.Config{},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("database connected")

	// =========================
	// DATABASE MIGRATION
	// =========================

	if err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},

		&models.Restaurant{},
		&models.Menu{},
		&models.RestaurantTable{},
		&models.Reservation{},

		&models.Order{},
		&models.OrderItem{},

		&models.Review{},
		&models.Reservation{},
		&models.Review{},
	); err != nil {
		log.Fatal(err)
	}

	log.Println("database migration completed")

	// =========================
	// REPOSITORIES
	// =========================

	userRepo := repository.NewUserRepository(db)

	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	restaurantRepo := repository.NewRestaurantRepository(db)

	menuRepo := repository.NewMenuRepository(db)
	tableRepo := repository.NewTableRepository(db)
	reservationRepo := repository.NewReservationRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	reviewRepo := repository.NewReviewRepository(db)


	// =========================
	// AUTH SERVICE
	// =========================

	authSvc := service.NewAuthService(
		userRepo,
		refreshTokenRepo,
		&cfg,
	)

	// =========================
	// AUTH HANDLER
	// =========================

	authHandler := handler.NewAuthHandler(
		authSvc,
	)

	// =========================
	// RESTAURANT SERVICE
	// =========================

	restaurantService := service.NewRestaurantService(
		restaurantRepo,
	)

	// =========================
	// RESTAURANT HANDLER
	// =========================

	restaurantHandler := handler.NewRestaurantHandler(
		restaurantService,
	)

	// =========================
	// MENU SERVICE
	// =========================

	menuService := service.NewMenuService(
		menuRepo,
		restaurantRepo,
	)
	
	tableService := service.NewTableService(
	tableRepo,
	restaurantRepo,
)

reservationService := service.NewReservationService(
	reservationRepo,
	restaurantRepo,
	tableRepo,
)

orderService := service.NewOrderService(
	orderRepo,
	menuRepo,
	restaurantRepo,
)

reviewService := service.NewReviewService(
	reviewRepo,
	restaurantRepo,
	orderRepo,
)
tableHandler := handler.NewTableHandler(
	tableService,
)
reservationHandler := handler.NewReservationHandler(
	reservationService,
)
orderHandler := handler.NewOrderHandler(
	orderService,
)
reviewHandler := handler.NewReviewHandler(
	reviewService,
)

	// =========================
	// MENU HANDLER
	// =========================

	menuHandler := handler.NewMenuHandler(
		menuService,
	)

	// =========================
	// APPLICATION
	// =========================

	app := &Application{
		Config:      &cfg,
		UserRepo:    userRepo,
		RefreshRepo: refreshTokenRepo,

		AuthService: authSvc,
		AuthHandler: authHandler,

		RestaurantService: restaurantService,
		RestaurantHandler: restaurantHandler,

		MenuService: menuService,
		MenuHandler: menuHandler,
		TableService: tableService,
	    TableHandler: tableHandler,
		ReservationService: reservationService,
		ReservationHandler: reservationHandler,
		OrderService: orderService,
		OrderHandler: orderHandler,
		ReviewService:      reviewService,
		ReviewHandler:      reviewHandler,
	}

	// =========================
	// SERVER
	// =========================

	log.Println("server listening on :8080")

	if err := http.ListenAndServe(
		":8080",
		app.routes(),
	); err != nil {
		log.Fatal(err)
	}
}