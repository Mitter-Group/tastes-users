package router

import (
	"github.com/chunnior/user-tastes-service/internal/app/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, loginHandler *handler.LoginHandler) {
	// app.Post("/login", handler.Login)
	app.Post("/login", loginHandler.Login)
	app.Post("/callback", loginHandler.Callback)
}
