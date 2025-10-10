package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	AppEnv        string
	Port          string
	PostgresUser  string
	PostgresPass  string
	PostgresDB    string
	PostgresHost  string
	PostgresPort  string
	PgSSLMode     string
	TelegramToken string
	JWTSecret     string
	MaxUsers      int
}

// LoadConfig for reading .env (или .env в dev)
func LoadConfig() (*Config, error) {
	// Загружаем .env only in dev
	if os.Getenv("APP_ENV") == "" {
		if err := godotenv.Load(".env"); err != nil {
			//or find in cmd/server
			_ = godotenv.Load("../.env")
		}
	}

	// Users limit
	maxUsers := 5
	if v := os.Getenv("MAX_USERS"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			maxUsers = parsed
		}
	}

	cfg := &Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		Port:          getEnv("PORT", "8080"),
		PostgresUser:  getEnv("POSTGRES_USER", "todo"),
		PostgresPass:  getEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:    getEnv("POSTGRES_DB", "todo_family"),
		PostgresHost:  getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:  getEnv("POSTGRES_PORT", "5432"),
		PgSSLMode:     getEnv("PGSSLMODE", "disable"),
		TelegramToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		MaxUsers:      maxUsers,
	}

	return cfg, nil
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.PostgresUser, c.PostgresPass, c.PostgresHost, c.PostgresPort, c.PostgresDB, c.PgSSLMode)
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
