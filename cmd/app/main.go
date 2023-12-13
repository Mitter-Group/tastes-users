package main

import (
	"log"

	"net/http"

	"github.com/chunnior/users/internal/app/handler"
	"github.com/chunnior/users/internal/app/router"
	"github.com/chunnior/users/internal/domain/login"
	"github.com/chunnior/users/internal/domain/provider"
	"github.com/chunnior/users/internal/infrastructure/logger"
	"github.com/chunnior/users/internal/repository/dynamodb"
	"github.com/chunnior/users/pkg/config"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load the configuration
	cfg := config.NewConfig()

	logger, err := logger.NewZapLogger()
	if err != nil {
		panic(err)
	}

	// Create a new Fiber instance
	app := fiber.New()

	userRepo := dynamodb.NewUserRepository(cfg, logger)

	loginService := login.NewLoginService(userRepo, cfg, &http.Client{}, logger)
	providerService := provider.NewProviderService(cfg, &http.Client{}, logger)

	loginHandler := handler.NewLoginHandler(*loginService)
	providerHandler := handler.NewProviderHandler(*providerService)

	// Setup the routes
	router.SetupRoutes(app, loginHandler, providerHandler)

	// Start the server
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
