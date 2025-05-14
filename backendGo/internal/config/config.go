package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
    Port     string
    DBConfig DBConfig
    JWT      JWTConfig
    Admin    AdminConfig
}

// DBConfig holds database configuration
type DBConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
    Secret     string
    Expiration time.Duration
}

// AdminConfig holds admin credentials
type AdminConfig struct {
    Username string
    Password string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found, using environment variables")
    }

    // Default JWT expiration to 7 days if not specified
    jwtExp := os.Getenv("JWT_EXPIRATION")
    if jwtExp == "" {
        jwtExp = "168h" // 7 days
    }

    expDuration, err := time.ParseDuration(jwtExp)
    if err != nil {
        log.Fatalf("Invalid JWT_EXPIRATION format: %v", err)
    }

    return &Config{
        Port: getEnv("PORT", "3000"),
        DBConfig: DBConfig{
            Host:     getEnv("POSTGRES_HOST", "localhost"),
            Port:     getEnv("POSTGRES_PORT", "5432"),
            User:     getEnv("POSTGRES_USER", "postgres"),
            Password: getEnv("POSTGRES_PASSWORD", "postgres"),
            DBName:   getEnv("POSTGRES_DB", "inventory"),
        },
        JWT: JWTConfig{
            Secret:     getEnv("JWT_SECRET", "default_jwt_secret"),
            Expiration: expDuration,
        },
        Admin: AdminConfig{
            Username: getEnv("ADMIN_USERNAME", "admin"),
            Password: getEnv("ADMIN_PASSWORD", "admin"),
        },
    }
}

// getEnv gets environment variable or returns fallback value
func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}