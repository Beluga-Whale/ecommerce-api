package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Beluga-Whale/ecommerce-api/config"
	"github.com/Beluga-Whale/ecommerce-api/internal/handlers"
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
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE",
		AllowHeaders: "Content-Type,Authorization",
		AllowCredentials: true,
	}))



	// NOTE - Create Repositories
	userRepo := repositories.NewUserRepository(config.DB)
	categoryRepo := repositories.NewCategoryRepository(config.DB)

	// NOTE - Utilities
	hashPassword := utils.NewPasswordUtil()
	jwtUtil := utils.NewJwt()

	// NOTE - Create Services
	userService := services.NewUserService(userRepo,hashPassword,jwtUtil)
	categoryService := services.NewCategoryService(categoryRepo)

	// NOTE - Create Handlers
	userHandler := handlers.NewUserHandler(userService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// NOTE - Set Up Routes
	routes.SetUpRoutes(app ,jwtUtil,userHandler,categoryHandler)
	
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