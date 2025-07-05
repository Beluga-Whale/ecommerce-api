package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type UserHandlerInterface interface{
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	GetProfile(c *fiber.Ctx) error
	UpdateProfile(c *fiber.Ctx) error
	// AddReview(c *fiber.Ctx) error
}


type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	// NOTE - Parse request body use DTO
	var req dto.RegisterRequestDTO
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
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:  req.Phone,
		BirthDate: req.BirthDate,
	}

	err := h.userService.Register(user)
	
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, err.Error())
	}

	return JSONSuccess(c, fiber.StatusCreated, "User registered successfully", nil)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	// NOTE - Parse request body use DTO
	var req dto.LoginRequestDTO
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

	token,userId,err := h.userService.Login(user)
	if err != nil {
		return JSONError(c, fiber.StatusUnauthorized, err.Error())
	}

	// NOTE - Set cookie
	c.Cookie(&fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour*72),
		Domain: ".belugaecommerce.xyz",
		HTTPOnly: true,
		Secure:true,
		SameSite: fiber.CookieSameSiteNoneMode, 
	})

	return JSONSuccess(c,fiber.StatusOK,"Login successful",dto.LoginResponseDTO{
		Token:  token,
		UserID: userId,
	})
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// NOTE - เอา UserIDจาก local
	// NOTE - ดึง userID จาก Locals แล้วแปลง string -> uint
	userIDStr, ok := c.Locals("userID").(string)

	if !ok {
		return JSONError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return JSONError(c, fiber.StatusInternalServerError, "Invalid user ID format")
	}

	user,err := h.userService.GetProfile(uint(userIDUint))

	return JSONSuccess(c,fiber.StatusOK,"Get Profile Successful",dto.UserProfileDTO{
		UserID: user.ID,
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Phone: user.Phone,
		BirthDate: user.BirthDate,
		Avatar: user.Avatar,
	})
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	var req dto.UserUpdateProfileDTO
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

	// NOTE - เอา UserIDจาก local
	// NOTE - ดึง userID จาก Locals แล้วแปลง string -> uint
	userIDStr, ok := c.Locals("userID").(string)

	if !ok {
		return JSONError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return JSONError(c, fiber.StatusInternalServerError, "Invalid user ID format")
	}

	err = h.userService.UpdateProfile(uint(userIDUint), req)
	
	if err != nil {
		return JSONError(c, fiber.StatusInternalServerError, "Failed to update profile")
	}

	return JSONSuccess(c, fiber.StatusOK, "Profile updated successfully", nil)
}

// func (h *UserHandler) AddReview(c *fiber.Ctx) error {
// 	var req dto.CreateReviewDTO
// 	if err := c.BodyParser(&req); err != nil {
// 		return JSONError(c, fiber.StatusBadRequest, "Invalid request body")
// 	}	

// 	// NOTE - Validate request body
// 	if err := Validate.Struct(req); err != nil {
// 		// NOTE - บอกว่า field ไหนผิด
// 		var messages []string
// 		for _, err := range err.(validator.ValidationErrors) {
// 			messages = append(messages, err.Field()+" is "+err.Tag())
// 		}
// 		return JSONError(c, fiber.StatusBadRequest, strings.Join(messages, ", "))
// 	}

// 	// NOTE - เอา UserIDจาก local
// 	// NOTE - ดึง userID จาก Locals แล้วแปลง string -> uint
// 	userIDStr, ok := c.Locals("userID").(string)

// 	if !ok {
// 		return JSONError(c, fiber.StatusUnauthorized, "Unauthorized")
// 	}

// 	userIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
// 	if err != nil {
// 		return JSONError(c, fiber.StatusInternalServerError, "Invalid user ID format")
// 	}
// }