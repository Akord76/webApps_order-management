package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// AppConfig holds every configuration value the web app needs.
// This app is an API client only: it never touches SQL Server directly,
// it talks to the order-management backend over HTTP.
type AppConfig struct {
	APIBaseURL string // e.g. http://localhost:8080/api

	ServerPort string

	JWTSecret string // must match the backend's secret so we can read claims from the cookie

	CookieName   string
	CookieSecure bool
}

var Cfg *AppConfig

func LoadConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("config: no .env file found, relying on OS environment variables")
	}

	Cfg = &AppConfig{
		APIBaseURL: getEnv("API_BASE_URL", "http://localhost:8083/api"),
		ServerPort: getEnv("SERVER_PORT", "8085"),
		JWTSecret:  mustGetEnv("JWT_SECRET"),

		CookieName:   getEnv("COOKIE_NAME", "token"),
		CookieSecure: getEnvBool("COOKIE_SECURE", false),
	}

	return Cfg
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		log.Fatalf("config: required env var %s is not set (check your .env file)", key)
	}
	return v
}

func getEnvBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
