package handlers

import (
	"strings"

	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	// NOTE - Parse request body use DTO
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return JSONError(c,fiber.StatusBadRequest, "Invalid request body")
	}

	// NOTE - Validate request body
		
	if err := Validate.Struct(req); err != nil {
		// NOTE - บอกว่า field ไหนผิด
		var messages []string
		for _, err := range err.(validator.ValidationErrors) {
			messages = append(messages, err.Field()+" is "+err.Tag())
		}
		return JSONError(c, fiber.StatusBadRequest, strings.Join(messages, ", "))
	}

	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.userService.Register(user)
	
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, err.Error())
	}

	return JSONSuccess(c, fiber.StatusCreated, "User registered successfully", nil)
}