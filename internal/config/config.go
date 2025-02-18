package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config structure
type Config struct {
	UserServiceURL  string
	AdminServiceURL string
	Port           string
}

// LoadConfig loads environment variables from a .env file
func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using defaults")
	}

	return &Config{
		UserServiceURL:  getEnv("USER_SERVICE_URL", "http://localhost:5001"),
		AdminServiceURL: getEnv("ADMIN_SERVICE_URL", "http://localhost:5003"),
		Port:           getEnv("GATEWAY_PORT", "8080"),
	}
}

// getEnv fetches an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
