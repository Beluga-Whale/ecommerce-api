package middleware

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// Middleware รับ jwtUtil เป็น dependency เพื่อให้ mock ได้ง่ายขึ้น
func AuthMiddleware(jwtUtil utils.JwtInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("jwt")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized - missing token",
			})
		}

		claims, err := jwtUtil.ParseJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized - invalid token",
			})
		}

		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)
		return c.Next()
	}
}


func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole")
		if userRole != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Forbidden - you don't have access",
			})
		}
		return c.Next()
	}
}
