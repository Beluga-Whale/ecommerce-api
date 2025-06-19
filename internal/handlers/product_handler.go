package handlers

import (
	"strconv"
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
		Title: 		 req.Title,
		Description: req.Description,
		IsFeatured:  req.IsFeatured,
		IsOnSale:    req.IsOnSale,
		SalePrice:   req.SalePrice,
		CategoryID:  req.CategoryID,
	}

	for _, v := range req.Variants{
		product.Variants = append(product.Variants, models.ProductVariant{
			Size: v.Size,
			Stock: v.Stock,
			SKU: v.SKU,
			Price: v.Price,
		})
	}

	for _, i := range req.Images {
		product.Images = append(product.Images, models.ProductImage{
			URL: i.URL,
		})
	}

	err := h.productService.CreateProduct(product)

	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, err.Error())
	}

	var variantsDTOs []dto.ProductVariantDTO

	for _, v:= range product.Variants {
		variantsDTOs = append(variantsDTOs, dto.ProductVariantDTO{
			Size: v.Size,
			Stock: v.Stock,
			SKU: v.SKU,
			Price: v.Price,
		})
	}

	var imageURLs []dto.ProductImageDTO
	for _, img := range product.Images {
		imageURLs = append(imageURLs, dto.ProductImageDTO{
			URL: img.URL,
		})
	}

	return JSONSuccess(c,fiber.StatusCreated, "Product created successfully", dto.ProductCreateResponseDTO{
		ID: 		 product.ID,
		Name:        product.Name,
		Title:       product.Title,
		Description: product.Description,
		Images:      imageURLs,
		IsFeatured:  product.IsFeatured,
		IsOnSale: 	 product.IsOnSale,
		SalePrice:   product.SalePrice,
		CategoryID:  product.CategoryID,
		Variants:    variantsDTOs,
	})
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
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
		Title: 		 req.Title,
		Description: req.Description,
		IsFeatured:  req.IsFeatured,
		IsOnSale:    req.IsOnSale,
		SalePrice:   req.SalePrice,
		CategoryID:  req.CategoryID,
	}

	for _, v := range req.Variants{
		product.Variants = append(product.Variants, models.ProductVariant{
			Size: v.Size,
			Stock: v.Stock,
			SKU: v.SKU,
			Price: v.Price,
		})
	}


	for _, i := range req.Images {
		product.Images = append(product.Images, models.ProductImage{
			URL: i.URL,
		})
	}


	err = h.productService.UpdateProduct(uint(id),product)
	if err != nil {
		return JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	var variantsDTOs []dto.ProductVariantDTO

	for _, v:= range product.Variants {
		variantsDTOs = append(variantsDTOs, dto.ProductVariantDTO{
			Size: v.Size,
			Stock: v.Stock,
			SKU: v.SKU,
			Price: v.Price,
		})
	}

	var imageURLs []dto.ProductImageDTO
	for _, img := range product.Images {
		imageURLs = append(imageURLs, dto.ProductImageDTO{
			URL: img.URL,
		})
	}

	return JSONSuccess(c, fiber.StatusOK, "Product updated successfully", dto.ProductUpdateResponseDTO{
		Name:        product.Name,
		Title: 		 product.Title,
		Description: product.Description,
		Images:      imageURLs,
		IsFeatured:  product.IsFeatured,
		IsOnSale:    product.IsOnSale,
		SalePrice:   product.SalePrice,
		CategoryID:  product.CategoryID,
		Variants:    variantsDTOs,
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
		return JSONError(c, fiber.StatusNotFound, "Product not found")
	}

	var variantsDTOs []dto.ProductVariantDTO
	for _, v:= range product.Variants {
		finalPrice := v.Price

		if product.IsOnSale && product.SalePrice != nil {
			finalPrice = v.Price - *product.SalePrice
			if finalPrice < 0 {
				finalPrice = 0 //NOTE - กันราคาติดลบ
			}
		}

		variantsDTOs = append(variantsDTOs, dto.ProductVariantDTO{
			Size: v.Size,
			Stock: v.Stock,
			SKU: v.SKU,
			Price: v.Price,
			FinalPrice: finalPrice,
		})
	}


	var imageURLs []dto.ProductImageDTO
	for _, img := range product.Images {
		imageURLs = append(imageURLs, dto.ProductImageDTO{
			URL: img.URL,
		})
	}


	return JSONSuccess(c, fiber.StatusOK, "Product retrieved successfully", dto.ProductUpdateResponseDTO{
		ID:          product.ID,
		Name:        product.Name,
		Title:       product.Title,
		Description: product.Description,
		Images:      imageURLs,
		IsFeatured:  product.IsFeatured,
		IsOnSale:    product.IsOnSale,
		SalePrice:   product.SalePrice,
		CategoryID:  product.CategoryID,
		Variants: 	 variantsDTOs,
	})
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	page := c.QueryInt("page",1)
	limit := c.QueryInt("limit",12)
	maxPrice := c.QueryInt("maxPrice",999999)
	minPrice := c.QueryInt("minPrice",0)
	searchName := c.Query("searchName","")
	category := c.Query("category","")
	size := c.Query("size","")

	categoryArr := strings.Split(category,",")
	sizeArr := strings.Split(size,",")

	var categoryIDs []int

	for _,i := range categoryArr {
		i = strings.TrimSpace(i)
		conId,err := strconv.Atoi(i)
		if i == "" {
			continue
		}	
		if err != nil {
			return JSONError(c, fiber.StatusInternalServerError, "can'n to convert categoryID string to int")
		}
		categoryIDs = append(categoryIDs, conId)
	}

	var sizeIDs []string
	for _,i := range sizeArr {
		i = strings.TrimSpace(i)
		if i == "" {
			continue
		}	
		
		sizeIDs = append(sizeIDs, i)
	}

	if limit <= 0 {
		limit = 1000000
	}

	if minPrice > maxPrice {
		return JSONError(c, fiber.StatusInternalServerError, "minPrice must be less than maxPrice")
	}

	products, pageTotal ,err := h.productService.GetAllProducts(uint(page),uint(limit),int64(minPrice),int64(maxPrice),searchName,categoryIDs,sizeIDs)

	if err != nil {
		return JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	var productsDTO []dto.ProductCreateResponseDTO

	for _,product := range products {
		var variantsDTOs []dto.ProductVariantDTO

		for _, v := range product.Variants {
			finalPrice := v.Price

			if product.IsOnSale && product.SalePrice != nil {
				finalPrice = v.Price - *product.SalePrice
				if finalPrice < 0 {
					finalPrice = 0
				}
			}

			variantsDTOs = append(variantsDTOs, dto.ProductVariantDTO{
				Size:       v.Size,
				Stock:      v.Stock,
				SKU:        v.SKU,
				Price:      v.Price,
				FinalPrice: finalPrice,
			})
		}

		var imageURLs []dto.ProductImageDTO
		for _, img := range product.Images {
			imageURLs = append(imageURLs, dto.ProductImageDTO{
				URL: img.URL,
			})
		}

		productsDTO = append(productsDTO, dto.ProductCreateResponseDTO{
			ID: 		 product.ID,	
			Name:        product.Name,
			Title:       product.Title,
			Description: product.Description,
			Images:      imageURLs,
			IsFeatured:  product.IsFeatured,
			IsOnSale:    product.IsOnSale,
			SalePrice:   product.SalePrice,
			CategoryID:  product.CategoryID,
			Variants:    variantsDTOs,
		})
	}

	return JSONSuccess(c, fiber.StatusOK, "Products retrieved successfully", fiber.Map{
		"products":productsDTO,
		"page": page,
		"limit" : limit,
		"pageTotal": pageTotal,
	})
}