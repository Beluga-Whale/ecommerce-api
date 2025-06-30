package handlers

import (
	"strconv"

	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

type ReviewHandlerInterface interface{
	GetUserReviews(c *fiber.Ctx) error 
	CreateReviews(c *fiber.Ctx) error 
}

type ReviewHandler struct {
	ReviewService services.ReviewServiceInterface
}

func NewReviewHandler(ReviewService services.ReviewServiceInterface) *ReviewHandler{
		return &ReviewHandler{ReviewService:ReviewService}
}

func (h *ReviewHandler) GetUserReviews(c *fiber.Ctx) error {
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
	reviews, err := h.ReviewService.GetReviewsByUserID(uint(userIDUint))
	if err != nil {
		return JSONError(c, fiber.StatusInternalServerError, "Failed to get reviews")
	}

	var reviewProduct []dto.ReviewResponse

	for _,item := range reviews {
		reviewProduct = append(reviewProduct,dto.ReviewResponse{
			ProductID: item.ProductID,
			Rating: item.Rating,
			Comment: item.Comment,
		} )
	}

	return JSONSuccess(c, fiber.StatusOK, "User reviews", reviewProduct)
}

func (h *ReviewHandler)	CreateReviews(c *fiber.Ctx) error {
	var req dto.CreateReviewDTO

	if err := c.BodyParser(&req); err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid request body")
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
	if err := h.ReviewService.CreateReview(uint(userIDUint), req); err != nil {
		return JSONError(c, fiber.StatusBadRequest, err.Error())
	}

	return JSONSuccess(c, fiber.StatusCreated, "Review created successfully", nil)
}