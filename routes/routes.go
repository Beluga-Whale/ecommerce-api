package routes

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/handlers"
	"github.com/Beluga-Whale/ecommerce-api/internal/middleware"
	"github.com/Beluga-Whale/ecommerce-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App, userHandler *handlers.UserHandler,jwtUtil utils.JwtInterface) {
	api := app.Group("/api")
	api.Post("/register", userHandler.Register)
	api.Post("/login",userHandler.Login)

	// Protected Routes
	protected := api.Group("/user", middleware.AuthMiddleware(jwtUtil))
	protected.Get("/profile", userHandler.GetProfile) // <- คุณต้องมี handler นี้ก่อน

	// Admin Routes
	admin := api.Group("/admin", middleware.AuthMiddleware(jwtUtil), middleware.RequireRole("admin"))
	admin.Get("/dashboard", userHandler.AdminDashboard) // <- handler นี้เฉพาะ admin
}