package handlers

import (
	"strconv"
	"strings"

	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
)

type CategoryHandler struct {
	categoryService services.CategoryServiceInterface
}

func NewCategoryHandler(categoryService services.CategoryServiceInterface) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetAllCategories()

	if err != nil {
		return JSONError(c, fiber.StatusInternalServerError, "Failed to fetch categories")
	}

	return JSONSuccess(c, fiber.StatusOK,"Get all categories successfully", categories)
}

func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	// NOTE - Parse request body use DTO
	var req dto.CategoryCreateDTO
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
		return JSONError(c, fiber.StatusBadRequest, "Validation error: "+strings.Join(messages, ", "))
	}

	category := &models.Category{
		Name: 	  req.Name,
	}

	err := h.categoryService.CreateCategory(category)
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, err.Error())
	}

	return JSONSuccess(c, fiber.StatusCreated, "Category created successfully", dto.CategoryCreateResponseDTO{
		Name: category.Name,
		Slug: category.Slug,
	})

}

func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	// NOTE - Get category ID from URL
	categoryID, _ := strconv.Atoi(c.Params("id")) 

	// NOTE - Parse request body use DTO
	var req dto.UpdateCategoryDTO
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
		return JSONError(c, fiber.StatusBadRequest, "Validation error: "+strings.Join(messages, ", "))
	}

	category := &models.Category{
		Name: 	  req.Name,
	}

	err := h.categoryService.UpdateCategory(uint(categoryID),category)
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, err.Error())
	}
	return JSONSuccess(c, fiber.StatusOK, "Category update successfully", dto.UpdateCategoryResponseDTO{
		Name: category.Name,
		Slug: slug.Make(category.Name) ,
	})
}

func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	// NOTE - Get category ID from URL
	categoryID, _ := strconv.Atoi(c.Params("id")) 

	// NOTE - Delete category
	err := h.categoryService.DeleteCategory(uint(categoryID))
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, err.Error())
	}

	return JSONSuccess(c, fiber.StatusOK, "Category deleted successfully", nil)
}