package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort                string
	SpotifyServiceURL      string
	AwsProfile             string
	AwsDynamoUserTableName string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		AppPort:                getEnv("PORT", "3000"),
		SpotifyServiceURL:      getEnv("SPOTIFY_SERVICE_URL", "http://localhost:8083"),
		AwsProfile:             getEnv("AWS_PROFILE", "default"),
		AwsDynamoUserTableName: getEnv("AWS_DYNAMO_USER_TABLE_NAME", "Users"),
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
