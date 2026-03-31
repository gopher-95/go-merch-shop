package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_USER     string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	DB_SSLMODE  string
}

// Строка для подключения к бд ...
func (cfg *Config) DatabaseURLString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME, cfg.DB_SSLMODE)
}

// Загрузка конфига ...
func LoadConf() *Config {
	_ = godotenv.Load()

	config := &Config{
		DB_USER:     getEnv("DB_USER", "postgres"),
		DB_PASSWORD: getEnv("DB_PASSWORD", ""),
		DB_HOST:     getEnv("DB_HOST", "localhost"),
		DB_PORT:     getEnv("DB_PORT", "5432"),
		DB_NAME:     getEnv("DB_NAME", "merch-shop"),
		DB_SSLMODE:  getEnv("DB_SSLMODE", "disable"),
	}

	if config.DB_PASSWORD == "" {
		log.Fatal("password is incorrect")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
