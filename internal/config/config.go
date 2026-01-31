package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	JWTSecret           string
	AutoCompleteMinutes int
}

func Load() *Config {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using system environment variables")
		}
	}
	minutes, err := strconv.Atoi(os.Getenv("AUTO_COMPLETE_MINUTES"))
	if err != nil || minutes <= 0 {
		log.Println("AUTO_COMPLETE_MINUTES not set or invalid, defaulting to 5")
		minutes = 5
	}

	cfg := &Config{
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),

		JWTSecret:           os.Getenv("JWT_SECRET"),
		AutoCompleteMinutes: minutes,
	}
	return cfg
}
