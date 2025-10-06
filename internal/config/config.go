package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
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
	MaxUploadSize int
}

func LoadConfig() (*Config, error) {
	//upload .env just for dev
	if os.Gotenv("APP_ENV") == "" {
		_ = godotenv.Load()
	}
}

maxUsers := 5
if v :=os.Getenv("MAX_USERS"); v != "" {
	if parsed, err := strconv.Atoi(v); err == nil {
		maxUsers = parsed
	}
}

cfg := &Config{
	AppEnv:        getenv("APP_ENV", "development"),
	Port:          getenv("PORT", ":8080"),
	PostgresUser:  getenv("POSTGRES_USER", "todo"),
	PostgresPass:  getenv("POSTGRES_PASS", ""),
	PostgresDB:    getenv("POSTGRES_DB", "todo-family"),
	PostgresHost:  getenv("POSTGRES_HOST", "localhost"),
	PostgresPort:  getenv("POSTGRES_PORT", "5432"),
	PgSSLMode:     getenv("POSTGRES_SSLMODE", "disable"),
	TelegramToken: getenv("TELEGRAM_TOKEN", ""),
	JWTSecret:     getenv("JWT_SECRET", ""),
	MaxUsers:      maxUsers,
}
return cfg, nil
}

func (c *Config) PostgreesDSn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.PostgresUser, c.PostgresPass, c.PostgresHost, c.PostgresPort, c.PostgresDB, c.PgSSLMode)
}

func getenv(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
