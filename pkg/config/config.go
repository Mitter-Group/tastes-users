package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort           string
	SpotifyServiceURL string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		AppPort:           getEnv("PORT", "3000"),
		SpotifyServiceURL: getEnv("SPOTIFY_SERVICE_URL", "http://localhost:8083"),
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
