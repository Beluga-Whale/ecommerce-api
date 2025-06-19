package routes

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/handlers"
	"github.com/Beluga-Whale/ecommerce-api/internal/middleware"
	"github.com/Beluga-Whale/ecommerce-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App, jwtUtil utils.JwtInterface, userHandler *handlers.UserHandler, categoryHandler *handlers.CategoryHandler, productHandler *handlers.ProductHandler, orderHandler *handlers.OrderHandler ) {
	api := app.Group("/api")
	api.Post("/register", userHandler.Register)
	api.Post("/login",userHandler.Login)
	api.Get("/category", categoryHandler.GetAll)
	api.Get("/product", productHandler.GetAllProducts)
	api.Get("/product/:id", productHandler.GetProductByID)
	// NOTE - Category Routes
	protectedCategoryAdmin := api.Group("/category", middleware.AuthMiddleware(jwtUtil),middleware.RequireRole("admin"))
	protectedCategoryAdmin.Post("/", categoryHandler.Create) 
	protectedCategoryAdmin.Put("/:id", categoryHandler.Update) 
	protectedCategoryAdmin.Delete("/:id", categoryHandler.Delete)

	// NOTE - Product Routes
	protectedProductAdmin := api.Group("/product", middleware.AuthMiddleware(jwtUtil), middleware.RequireRole("admin"))
	protectedProductAdmin.Post("/", productHandler.CreateProduct) 
	protectedProductAdmin.Put("/:id", productHandler.UpdateProduct)
	protectedProductAdmin.Delete("/:id", productHandler.DeleteProduct)

	// NOTE - Order Route
	// NOTE - User use createOnly
	protectedOrderUser := api.Group("/user/order", middleware.AuthMiddleware(jwtUtil), middleware.RequireRole("user"))
	protectedOrderUser.Post("/", orderHandler.CreateOrder)


	// Protected Routes
	protected := api.Group("/user", middleware.AuthMiddleware(jwtUtil))
	protected.Get("/profile", userHandler.GetProfile) // <- คุณต้องมี handler นี้ก่อน

	// Admin Routes
	admin := api.Group("/admin", middleware.AuthMiddleware(jwtUtil), middleware.RequireRole("admin"))
	admin.Get("/dashboard", userHandler.AdminDashboard) // <- handler นี้เฉพาะ admin
}