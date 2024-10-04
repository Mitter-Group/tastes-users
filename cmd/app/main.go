package main

import (
	"log"

	"net/http"

	"github.com/chunnior/users/internal/app/handler"
	"github.com/chunnior/users/internal/app/router"
	"github.com/chunnior/users/internal/domain/login"
	"github.com/chunnior/users/internal/domain/provider"
	"github.com/chunnior/users/internal/infrastructure/aws/sqs"
	"github.com/chunnior/users/internal/infrastructure/logger"
	"github.com/chunnior/users/internal/repository/dynamodb"
	"github.com/chunnior/users/pkg/config"

	"github.com/gofiber/fiber/v2"
)

func main() {

	logger, err := logger.NewZapLogger()
	if err != nil {
		panic(err)
	}
	// Load the configuration
	cfg := config.NewConfig()

	// secretKey := cfg.SecretKey
	// encryptedAPIKey := cfg.EncryptedAPIKey
	// Create a new Fiber instance
	app := fiber.New()

	// authMiddleware := middleware.AuthMiddleware(encryptedAPIKey, secretKey)

	// app.Use(authMiddleware)

	userRepo := dynamodb.NewUserRepository(cfg, logger)

	sqsClient, err := sqs.NewSQS(cfg)
	if err != nil {
		panic(err)
	}

	loginService := login.NewLoginService(userRepo, sqsClient, cfg, &http.Client{}, logger)
	providerService := provider.NewProviderService(cfg, &http.Client{}, logger)

	healthHandler := handler.NewHealthHandler()
	loginHandler := handler.NewLoginHandler(*loginService)
	providerHandler := handler.NewProviderHandler(*providerService)

	// Setup the routes
	router.SetupRoutes(app, loginHandler, providerHandler, healthHandler)

	// Start the server
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
