package handlers

import (
	"strings"
	"time"

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

func (h *UserHandler) Login(c *fiber.Ctx) error {
	// NOTE - Parse request body use DTO
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid request body")
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

	token,err := h.userService.Login(user)
	if err != nil {
		return JSONError(c, fiber.StatusUnauthorized, err.Error())
	}

	// NOTE - Set cookie
	c.Cookie(&fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour*72),
		// Domain: ".belugatasks.dev",
		HTTPOnly: true,
		Secure:false,
		SameSite: fiber.CookieSameSiteNoneMode, 
		
	})

	return JSONSuccess(c,fiber.StatusOK,"Login successful",dto.LoginResponse{
		Token:  token,
	})
}