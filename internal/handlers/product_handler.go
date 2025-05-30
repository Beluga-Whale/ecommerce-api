package handlers

import (
	"strings"

	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type ProductHandlerInterface interface{
	CreateProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
	GetProductByID(c *fiber.Ctx) error
	GetAllProducts(c *fiber.Ctx) error
}

type ProductHandler struct {
	productService services.ProductServiceInterface
}

func NewProductHandler(productService services.ProductServiceInterface) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	// NOTE - Parse request body use DTO
	var req dto.ProductCreateDTO
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
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Image:       req.Image,
		Stock:       req.Stock,
		IsFeatured: req.IsFeatured,
		IsOnSale: 	req.IsOnSale,
		CategoryID: req.CategoryID,
	}

	err := h.productService.CreateProduct(product)

	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, err.Error())
	}

	return JSONSuccess(c,fiber.StatusCreated, "Product created successfully", dto.ProductCreateResponseDTO{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Image:       product.Image,
		Stock:       product.Stock,
		IsFeatured: product.IsFeatured,
		IsOnSale: 	product.IsOnSale,
		CategoryID: product.CategoryID,
	})
}

func (h* ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	var req dto.ProductUpdateDTO
	if err := c.BodyParser(&req); err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := Validate.Struct(req); err != nil {
		var messages []string
		for _, err := range err.(validator.ValidationErrors) {
			messages = append(messages, err.Field()+" is "+err.Tag())
		}
		return JSONError(c, fiber.StatusBadRequest, "Validation error: "+strings.Join(messages, ", "))
	}

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Image:       req.Image,
		Stock:       req.Stock,
		IsFeatured:  req.IsFeatured,
		IsOnSale:    req.IsOnSale,
		CategoryID:  req.CategoryID,
	}

	err = h.productService.UpdateProduct(uint(id),product)
	if err != nil {
		return JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return JSONSuccess(c, fiber.StatusOK, "Product updated successfully", dto.ProductUpdateResponseDTO{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Image:       product.Image,
		Stock:       product.Stock,
		IsFeatured:  product.IsFeatured,
		IsOnSale:    product.IsOnSale,
		CategoryID:  product.CategoryID,
	})
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	err = h.productService.DeleteProduct(uint(id))
	if err != nil {
		return JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return JSONSuccess(c, fiber.StatusOK, "Product deleted successfully", nil)
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		return JSONError(c, fiber.StatusInternalServerError,  err.Error())
	}
	if product == nil {
		return JSONError(c, fiber.StatusNotFound, err.Error())
	}

	return JSONSuccess(c, fiber.StatusOK, "Product retrieved successfully", dto.ProductUpdateResponseDTO{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Image:       product.Image,
		Stock:       product.Stock,
		IsFeatured:  product.IsFeatured,
		IsOnSale:    product.IsOnSale,
		CategoryID:  product.CategoryID,
	})
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		return JSONError(c, fiber.StatusInternalServerError, err.Error())
	}
	return JSONSuccess(c, fiber.StatusOK, "Products retrieved successfully", products)
}