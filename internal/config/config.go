package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
	ServerPort string
}

// Строка для подключения к бд
func (cfg *Config) DatabaseURLString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)
}

// Загрузка конфига
func LoadConf() *Config {
	_ = godotenv.Load()

	config := &Config{
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "merch-shop"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	if config.DBPassword == "" {
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
