package handlers_test

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Beluga-Whale/ecommerce-api/internal/handlers"
	services "github.com/Beluga-Whale/ecommerce-api/internal/services/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProduct(t *testing.T) {
	t.Run("Create Product",func(t *testing.T) {
		productService := services.NewProductServiceMock()

		productHandler := handlers.NewProductHandler(productService)

		app := fiber.New()
		app.Post("/product",productHandler.CreateProduct)

		productService.On("CreateProduct",mock.Anything).Return(nil)

		reqBody:= []byte(`{
			"name": "T-Shirt",
			"title": "test title",
			"description": "<p>TEST</p>",
			"images": [
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060281/Beluga_Whale_gtyr3j.webp"
				},
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060282/beluga-whale-swimming-Norway_oelbul.jpg"
				},
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060282/03_white-tee_model-front-scaled-scaled_i69es8.jpg"
				}
			],
			"isFeatured": false,
			"isOnSale": true,
			"salePrice": 20,
			"categoryId": 9,
			"variants": [
				{
					"size": "S",
					"stock": 3,
					"sku": "T-SHIRT-S",
					"price": 130
				},
				{
					"size": "M",
					"stock": 1,
					"sku": "T-SHIRT-M",
					"price": 140
				}
			]
		}`)

		req :=httptest.NewRequest("POST","/product",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Product created successfully")
	})
	t.Run("Request Nil",func(t *testing.T) {
		productService := services.NewProductServiceMock()

		productHandler := handlers.NewProductHandler(productService)

		app := fiber.New()
		app.Post("/product",productHandler.CreateProduct)


		req :=httptest.NewRequest("POST","/product",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request body")
	})
	t.Run("Validate Error",func(t *testing.T) {
		productService := services.NewProductServiceMock()

		productHandler := handlers.NewProductHandler(productService)

		app := fiber.New()
		app.Post("/product",productHandler.CreateProduct)
		
		reqBody:= []byte(`{
			"title": "test title",
			"description": "<p>TEST</p>",
			"images": [
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060281/Beluga_Whale_gtyr3j.webp"
				},
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060282/beluga-whale-swimming-Norway_oelbul.jpg"
				},
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060282/03_white-tee_model-front-scaled-scaled_i69es8.jpg"
				}
			],
			"isFeatured": false,
			"isOnSale": true,
			"salePrice": 20,
			"categoryId": 9,
			"variants": [
				{
					"size": "S",
					"stock": 3,
					"sku": "T-SHIRT-S",
					"price": 130
				},
				{
					"size": "M",
					"stock": 1,
					"sku": "T-SHIRT-M",
					"price": 140
				}
			]
		}`)

		req :=httptest.NewRequest("POST","/product",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Name is required")
	})
	t.Run("Error to create Product",func(t *testing.T) {
		productService := services.NewProductServiceMock()

		productHandler := handlers.NewProductHandler(productService)

		app := fiber.New()
		app.Post("/product",productHandler.CreateProduct)

		productService.On("CreateProduct",mock.Anything).Return(errors.New("Error to create product"))

		reqBody:= []byte(`{
			"name": "T-Shirt",
			"title": "test title",
			"description": "<p>TEST</p>",
			"images": [
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060281/Beluga_Whale_gtyr3j.webp"
				},
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060282/beluga-whale-swimming-Norway_oelbul.jpg"
				},
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060282/03_white-tee_model-front-scaled-scaled_i69es8.jpg"
				}
			],
			"isFeatured": false,
			"isOnSale": true,
			"salePrice": 20,
			"categoryId": 9,
			"variants": [
				{
					"size": "S",
					"stock": 3,
					"sku": "T-SHIRT-S",
					"price": 130
				},
				{
					"size": "M",
					"stock": 1,
					"sku": "T-SHIRT-M",
					"price": 140
				}
			]
		}`)

		req :=httptest.NewRequest("POST","/product",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Error to create product")
	})
}

func TestUpdateProduct(t *testing.T) {
	t.Run("Update Success",func(t *testing.T) {
		productService := services.NewProductServiceMock()

		productHandler := handlers.NewProductHandler(productService)

		app := fiber.New()
		app.Put("/product/:id",productHandler.UpdateProduct)

		productService.On("UpdateProduct",mock.Anything,mock.Anything).Return(nil)

		reqBody:= []byte(`{
			"name": "T-Shirt",
			"title": "test title",
			"description": "<p>TEST</p>",
			"images": [
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060281/Beluga_Whale_gtyr3j.webp"
				},
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060282/beluga-whale-swimming-Norway_oelbul.jpg"
				},
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060282/03_white-tee_model-front-scaled-scaled_i69es8.jpg"
				}
			],
			"isFeatured": false,
			"isOnSale": true,
			"salePrice": 20,
			"categoryId": 9,
			"variants": [
				{
					"size": "S",
					"stock": 3,
					"sku": "T-SHIRT-S",
					"price": 130
				},
				{
					"size": "M",
					"stock": 1,
					"sku": "T-SHIRT-M",
					"price": 140
				}
			]
		}`)

		req :=httptest.NewRequest("PUT","/product/1",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Product updated successfully")
	})
	t.Run("Invalid Param id",func(t *testing.T) {
		productService := services.NewProductServiceMock()

		productHandler := handlers.NewProductHandler(productService)

		app := fiber.New()
		app.Put("/product/:id",productHandler.UpdateProduct)

		productService.On("UpdateProduct",mock.Anything,mock.Anything).Return(nil)

		reqBody:= []byte(`{
			"name": "T-Shirt",
			"title": "test title",
			"description": "<p>TEST</p>",
			"images": [
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060281/Beluga_Whale_gtyr3j.webp"
				},
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060282/beluga-whale-swimming-Norway_oelbul.jpg"
				},
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060282/03_white-tee_model-front-scaled-scaled_i69es8.jpg"
				}
			],
			"isFeatured": false,
			"isOnSale": true,
			"salePrice": 20,
			"categoryId": 9,
			"variants": [
				{
					"size": "S",
					"stock": 3,
					"sku": "T-SHIRT-S",
					"price": 130
				},
				{
					"size": "M",
					"stock": 1,
					"sku": "T-SHIRT-M",
					"price": 140
				}
			]
		}`)

		req :=httptest.NewRequest("PUT","/product/failId",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid product ID")
	})

	t.Run("Validate Request Body",func(t *testing.T) {
		productService := services.NewProductServiceMock()

		productHandler := handlers.NewProductHandler(productService)

		app := fiber.New()
		app.Put("/product/:id",productHandler.UpdateProduct)

		productService.On("UpdateProduct",mock.Anything,mock.Anything).Return(nil)

		req :=httptest.NewRequest("PUT","/product/1",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request body")
	})
	t.Run("Update Success",func(t *testing.T) {
		productService := services.NewProductServiceMock()

		productHandler := handlers.NewProductHandler(productService)

		app := fiber.New()
		app.Put("/product/:id",productHandler.UpdateProduct)

		productService.On("UpdateProduct",mock.Anything,mock.Anything).Return(nil)

		reqBody:= []byte(`{
			"title": "test title",
			"description": "<p>TEST</p>",
			"images": [
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060281/Beluga_Whale_gtyr3j.webp"
				},
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060282/beluga-whale-swimming-Norway_oelbul.jpg"
				},
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060282/03_white-tee_model-front-scaled-scaled_i69es8.jpg"
				}
			],
			"isFeatured": false,
			"isOnSale": true,
			"salePrice": 20,
			"categoryId": 9,
			"variants": [
				{
					"size": "S",
					"stock": 3,
					"sku": "T-SHIRT-S",
					"price": 130
				},
				{
					"size": "M",
					"stock": 1,
					"sku": "T-SHIRT-M",
					"price": 140
				}
			]
		}`)

		req :=httptest.NewRequest("PUT","/product/1",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Name is required")
	})
	t.Run("Error UpdateProduct",func(t *testing.T) {
		productService := services.NewProductServiceMock()

		productHandler := handlers.NewProductHandler(productService)

		app := fiber.New()
		app.Put("/product/:id",productHandler.UpdateProduct)

		productService.On("UpdateProduct",mock.Anything,mock.Anything).Return(errors.New("Error update product"))

		reqBody:= []byte(`{
			"name": "T-Shirt",
			"title": "test title",
			"description": "<p>TEST</p>",
			"images": [
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060281/Beluga_Whale_gtyr3j.webp"
				},
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060282/beluga-whale-swimming-Norway_oelbul.jpg"
				},
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060282/03_white-tee_model-front-scaled-scaled_i69es8.jpg"
				}
			],
			"isFeatured": false,
			"isOnSale": true,
			"salePrice": 20,
			"categoryId": 9,
			"variants": [
				{
					"size": "S",
					"stock": 3,
					"sku": "T-SHIRT-S",
					"price": 130
				},
				{
					"size": "M",
					"stock": 1,
					"sku": "T-SHIRT-M",
					"price": 140
				}
			]
		}`)

		req :=httptest.NewRequest("PUT","/product/1",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Error update product")
	})
}

func TestDeleteProduct( t *testing.T) {
	t.Run("Delete Success",func(t *testing.T) {
		productService := services.NewProductServiceMock()

		productHandler := handlers.NewProductHandler(productService)

		app := fiber.New()
		app.Delete("/product/:id",productHandler.DeleteProduct)

		productService.On("DeleteProduct",mock.Anything).Return(nil)


		req :=httptest.NewRequest("DELETE","/product/1",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Product deleted successfully")	
	})

	t.Run("Invalid Param ID",func(t *testing.T) {
		productService := services.NewProductServiceMock()

		productHandler := handlers.NewProductHandler(productService)

		app := fiber.New()
		app.Delete("/product/:id",productHandler.DeleteProduct)

		productService.On("DeleteProduct",mock.Anything).Return(nil)


		req :=httptest.NewRequest("DELETE","/product/fail",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid product ID")	
	})
	t.Run("Error to Delete",func(t *testing.T) {
		productService := services.NewProductServiceMock()

		productHandler := handlers.NewProductHandler(productService)

		app := fiber.New()
		app.Delete("/product/:id",productHandler.DeleteProduct)

		productService.On("DeleteProduct",mock.Anything).Return(errors.New("Error to delete product"))


		req :=httptest.NewRequest("DELETE","/product/1",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Error to delete product")	
	})
}