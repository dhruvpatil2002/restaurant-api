package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
 "gorm.io/gorm/logger"       
 
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

	RestaurantService  *service.RestaurantService
	RestaurantHandler  *handler.RestaurantHandler
	MenuService        *service.MenuService
	MenuHandler        *handler.MenuHandler
	TableService       *service.TableService
	TableHandler       *handler.TableHandler
	ReservationService *service.ReservationService
	ReservationHandler *handler.ReservationHandler
	OrderService       *service.OrderService
	OrderHandler       *handler.OrderHandler
	ReviewService      *service.ReviewService
	ReviewHandler      *handler.ReviewHandler
}

func main() {
	
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}


	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	
	db, err := gorm.Open(
		postgres.Open(cfg.DBURL),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("database connected")

	
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
	); err != nil {
		log.Fatal(err)
	}

	log.Println("database migration completed")

	

	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	restaurantRepo := repository.NewRestaurantRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	tableRepo := repository.NewTableRepository(db)
	reservationRepo := repository.NewReservationRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	

	authSvc := service.NewAuthService(
		userRepo,
		refreshTokenRepo,
		&cfg,
	)

	restaurantService := service.NewRestaurantService(restaurantRepo)
	menuService := service.NewMenuService(menuRepo, restaurantRepo)
	tableService := service.NewTableService(tableRepo, restaurantRepo)
	reservationService := service.NewReservationService(reservationRepo, restaurantRepo, tableRepo)
	orderService := service.NewOrderService(orderRepo, menuRepo, restaurantRepo)
	reviewService := service.NewReviewService(reviewRepo, restaurantRepo, orderRepo)

	

	authHandler := handler.NewAuthHandler(authSvc)
	restaurantHandler := handler.NewRestaurantHandler(restaurantService)
	menuHandler := handler.NewMenuHandler(menuService)
	tableHandler := handler.NewTableHandler(tableService)
	reservationHandler := handler.NewReservationHandler(reservationService)
	orderHandler := handler.NewOrderHandler(orderService)
	reviewHandler := handler.NewReviewHandler(reviewService)

	

	app := &Application{
		Config:             &cfg,
		UserRepo:           userRepo,
		RefreshRepo:        refreshTokenRepo,
		AuthService:        authSvc,
		AuthHandler:        authHandler,
		RestaurantService:  restaurantService,
		RestaurantHandler:  restaurantHandler,
		MenuService:        menuService,
		MenuHandler:        menuHandler,
		TableService:       tableService,
		TableHandler:       tableHandler,
		ReservationService: reservationService,
		ReservationHandler: reservationHandler,
		OrderService:       orderService,
		OrderHandler:       orderHandler,
		ReviewService:      reviewService,
		ReviewHandler:      reviewHandler,
	}






	

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      app.routes(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	
	go func() {
		log.Println("server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped")
}