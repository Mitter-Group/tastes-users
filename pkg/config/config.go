package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort                string
	SpotifyServiceURL      string
	YoutubeServiceURL      string
	AwsProfile             string
	AwsDynamoUserTableName string
	EncryptedAPIKey        string
	SecretKey              string
	NewUserQueueURL        string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		AppPort:                getEnv("PORT", "3000"),
		SpotifyServiceURL:      getEnv("SPOTIFY_SERVICE_URL", "http://localhost:8083"),
		YoutubeServiceURL:      getEnv("YOUTUBE_SERVICE_URL", "http://localhost:8085"),
		AwsProfile:             getEnv("AWS_PROFILE", "default"),
		AwsDynamoUserTableName: getEnv("AWS_DYNAMO_USER_TABLE_NAME", "Users"),
		EncryptedAPIKey:        getEnv("ENCRYPTED_API_KEY", ""),
		SecretKey:              getEnv("SECRET_KEY", ""),
		NewUserQueueURL:        getEnv("QUEUE_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
