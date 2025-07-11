package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Beluga-Whale/ecommerce-api/config"
	"github.com/Beluga-Whale/ecommerce-api/internal/handlers"
	"github.com/Beluga-Whale/ecommerce-api/internal/jobs"
	"github.com/Beluga-Whale/ecommerce-api/internal/repositories"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/Beluga-Whale/ecommerce-api/internal/utils"
	"github.com/Beluga-Whale/ecommerce-api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// NOTE - LoadEnv
	config.LoadEnv()

	// NOTE - Connect DB
	config.ConnectDB()

	// NOTE - Fiber
	app := fiber.New()

	// NOTE - Use cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, http://localhost:3001,https://belugaecommerce.xyz",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE",
		AllowHeaders: "Content-Type,Authorization",
		AllowCredentials: true,
	}))

	// NOTE - Create Repositories
	userRepo := repositories.NewUserRepository(config.DB)
	categoryRepo := repositories.NewCategoryRepository(config.DB)
	productRepo := repositories.NewProductRepository(config.DB)
	orderRepo := repositories.NewOrderRepository(config.DB)
	reviewRepo := repositories.NewReviewRepository(config.DB)

	// NOTE - Utilities
	hashPassword := utils.NewPasswordUtil()
	jwtUtil := utils.NewJwt()
	productUtil := utils.NewProductUtil()

	// NOTE - Create Services
	userService := services.NewUserService(userRepo,hashPassword,jwtUtil)
	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo,categoryRepo)
	orderService := services.NewOrderService(config.DB,orderRepo, productUtil)
	reviewService := services.NewReviewService(reviewRepo)
	
	// NOTE - Create Handlers
	userHandler := handlers.NewUserHandler(userService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)
	PaymentHandler := handlers.NewStripeHandler(orderService)
	ReviewHandler := handlers.NewReviewHandler(reviewService)

	// NOTE - Set Up Routes
	routes.SetUpRoutes(app ,jwtUtil,userHandler,categoryHandler,productHandler,orderHandler,PaymentHandler,ReviewHandler)
	

	// NOTE -ทำงานเพื่อการนับถอยหลังเช็ค order
	jobs.StartOrderExpirationJob(config.DB, orderService)

	port := os.Getenv("PORT_API")

	if port =="" {
		port =":8080"
	}

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Server running on port %s\n", port)
	
    // NOTE -เช็ค error จาก Listen
    if err := app.Listen(addr); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}