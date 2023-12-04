package main

import (
	"log"

	"net/http"

	"github.com/chunnior/user-tastes-service/internal/app/handler"
	"github.com/chunnior/user-tastes-service/internal/app/router"
	"github.com/chunnior/user-tastes-service/internal/domain/login"
	"github.com/chunnior/user-tastes-service/pkg/config"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load the configuration
	cfg := config.NewConfig()

	// Create a new Fiber instance
	app := fiber.New()

	loginService := login.NewLoginService(cfg, &http.Client{})

	loginHandler := handler.NewLoginHandler(*loginService)

	// Setup the routes
	router.SetupRoutes(app, loginHandler)

	// Start the server
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
