package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// App
	Port string

	// JWT
	JWTAccessSecret      string
	JWTRefreshSecret     string
	JWTAccessExpiration  time.Duration
	JWTRefreshExpiration time.Duration

	// OAuth Yandex
	YandexClientID     string
	YandexClientSecret string
	YandexRedirectURL  string

	// OAuth VK
	VKClientID     string
	VKClientSecret string
	VKRedirectURL  string
}

func Load() (*Config, error) {
	godotenv.Load()

	accessExp, _ := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRATION", "15m"))
	refreshExp, _ := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRATION", "168h"))

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "student"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "paving_tiles_db"),
		Port:       getEnv("PORT", "4200"),

		JWTAccessSecret:      getEnv("JWT_ACCESS_SECRET", "access_secret"),
		JWTRefreshSecret:     getEnv("JWT_REFRESH_SECRET", "refresh_secret"),
		JWTAccessExpiration:  accessExp,
		JWTRefreshExpiration: refreshExp,

		YandexClientID:     getEnv("YANDEX_CLIENT_ID", ""),
		YandexClientSecret: getEnv("YANDEX_CLIENT_SECRET", ""),
		YandexRedirectURL:  getEnv("YANDEX_REDIRECT_URL", ""),

		VKClientID:     getEnv("VK_CLIENT_ID", ""),
		VKClientSecret: getEnv("VK_CLIENT_SECRET", ""),
		VKRedirectURL:  getEnv("VK_REDIRECT_URL", ""),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
