package router

import (
	"github.com/chunnior/users/internal/app/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, loginHandler *handler.LoginHandler, providerHandler *handler.ProviderHandler) {
	app.Post("/login", loginHandler.Login)
	app.Post("/callback", loginHandler.Callback)
	app.Get("/:provider/:dataType/:userId", providerHandler.ProviderInfo)
}
