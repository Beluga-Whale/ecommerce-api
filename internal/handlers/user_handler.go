package handlers

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return JSONError(c,fiber.StatusBadRequest, "Invalid request body")
	}

	err := h.userService.Register(user)
	
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, err.Error())
	}

	return JSONSuccess(c, fiber.StatusCreated, "User registered successfully", nil)
}