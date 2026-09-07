package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv             string
	HTTPAddr           string
	DatabaseURL        string
	RedisURL           string
	JWTSecret          string
	JWTAccessTTL       time.Duration
	JWTRefreshTTL      time.Duration
	CORSAllowedOrigins []string
	ShutdownTimeout    time.Duration
}

func Load() (Config, error) {
	accessMinutes, err := parseInt(getEnv("JWT_ACCESS_TTL_MINUTES", "15"), "JWT_ACCESS_TTL_MINUTES")
	if err != nil {
		return Config{}, err
	}
	refreshDays, err := parseInt(getEnv("JWT_REFRESH_TTL_DAYS", "7"), "JWT_REFRESH_TTL_DAYS")
	if err != nil {
		return Config{}, err
	}
	shutdownSeconds, err := parseInt(getEnv("SHUTDOWN_TIMEOUT_SECONDS", "10"), "SHUTDOWN_TIMEOUT_SECONDS")
	if err != nil {
		return Config{}, err
	}

	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if jwtSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	dbURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dbURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	cfg := Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		HTTPAddr:           getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL:        dbURL,
		RedisURL:           getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:          jwtSecret,
		JWTAccessTTL:       time.Duration(accessMinutes) * time.Minute,
		JWTRefreshTTL:      time.Duration(refreshDays) * 24 * time.Hour,
		CORSAllowedOrigins: splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:8080")),
		ShutdownTimeout:    time.Duration(shutdownSeconds) * time.Second,
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func parseInt(raw, name string) (int, error) {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", name, err)
	}
	return value, nil
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
