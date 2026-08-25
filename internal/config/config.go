package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv           string
	HTTPPort         int
	Database         DatabaseConfig
	DBURL            string
	JWTSigningKey    []byte
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	httpPort, err := getInt("HTTP_PORT", 8080)
	if err != nil {
		return Config{}, err
	}

	dbPort, err := getInt("DB_PORT", 5432)
	if err != nil {
		return Config{}, err
	}

	accessMins, err := getInt("JWT_ACCESS_EXPIRY_MINUTES", 20)
	if err != nil {
		return Config{}, err
	}

	refreshDays, err := getInt("JWT_REFRESH_EXPIRY_DAYS", 14)
	if err != nil {
		return Config{}, err
	}

	dbHost := getString("DB_HOST", "localhost")
	dbUser := getString("DB_USER", "restaurant_user")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := getString("DB_NAME", "restaurant_db")
	dbSSLMode := getString("DB_SSLMODE", "disable")

	dbURL := getString("DB_URL", "")
	if dbURL == "" {
		dbURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode,
		)
	}

	signingKeyStr := getString("JWT_SIGNING_KEY", "default-signing-key")

	config := Config{
		AppEnv:           getString("APP_ENV", "development"),
		HTTPPort:         httpPort,
		DBURL:            dbURL,
		JWTSigningKey:    []byte(signingKeyStr),
		JWTAccessExpiry:  time.Duration(accessMins) * time.Minute,
		JWTRefreshExpiry: time.Duration(refreshDays) * 24 * time.Hour,
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
			SSLMode:  dbSSLMode,
		},
	}

	if config.Database.Password == "" {
		return Config{}, fmt.Errorf("DB_PASSWORD is required")
	}

	return config, nil
}

func getString(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getInt(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer", key)
	}

	return result, nil
}