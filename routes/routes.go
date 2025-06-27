package routes

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/handlers"
	"github.com/Beluga-Whale/ecommerce-api/internal/middleware"
	"github.com/Beluga-Whale/ecommerce-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App, jwtUtil utils.JwtInterface, userHandler *handlers.UserHandler, categoryHandler *handlers.CategoryHandler, productHandler *handlers.ProductHandler, orderHandler *handlers.OrderHandler, paymentHandler *handlers.StripeHandler ) {


	api := app.Group("/api")
	api.Post("/register", userHandler.Register)
	api.Post("/login",userHandler.Login)
	api.Get("/category", categoryHandler.GetAll)
	api.Get("/product", productHandler.GetAllProducts)
	api.Get("/product/:id", productHandler.GetProductByID)
	api.Get("/user/order/:id", orderHandler.GetOrderByID)

	// NOTE  - Payment	
	api.Post("/stripe/payment-intent",paymentHandler.CreatePaymentIntent)
	api.Post("/stripe/webhook", paymentHandler.Webhook)

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
	protectedOrderUser.Patch("/", orderHandler.UpdateStatusOrder)
	protectedOrderUser.Get("/",orderHandler.GetAllOrderByUserId)
	protectedOrderUser.Patch("/:id/status",orderHandler.UpdateOrderStatusByUser)	
	
	
	// NOTE - Admin use Order
	protectedOrderAdmin := api.Group("/admin/order", middleware.AuthMiddleware(jwtUtil), middleware.RequireRole("admin"))
	protectedOrderAdmin.Get("/",orderHandler.GetAllOrders)
}