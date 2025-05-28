package routes

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App, userHandler *handlers.UserHandler) {
	api := app.Group("/api")
	api.Post("/register", userHandler.Register)
	api.Post("/login",userHandler.Login)
}