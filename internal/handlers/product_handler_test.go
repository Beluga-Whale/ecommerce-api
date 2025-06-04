package handlers_test

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Beluga-Whale/ecommerce-api/internal/handlers"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	services "github.com/Beluga-Whale/ecommerce-api/internal/services/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	t.Run("CreateProduct Success", func(t *testing.T) {
		product := &models.Product{
			Name:        "CPU",
			Description: "descriptionTest",
			Price:       100,
			Image:       "https://www.google.com",
			Stock:       10,
			IsFeatured:  false,
			IsOnSale:    false,
			CategoryID:  3,
		}

		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Post("/product", productHandler.CreateProduct)

		productService.On("CreateProduct",product).Return(nil)
		
		reqBody := []byte(`{"name":"CPU","description":"descriptionTest","Price":100,"image":"https://www.google.com","stock":10,"isFeatured":false,"isOnSale":false,"categoryID":3}`)

		req := httptest.NewRequest("POST","/product", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Product created successfully")
	})

	t.Run("Invalid request body", func(t *testing.T) {
		product := &models.Product{
			Name:        "CPU",
			Description: "descriptionTest",
			Price:       100,
			Image:       "https://www.google.com",
			Stock:       10,
			IsFeatured:  false,
			IsOnSale:    false,
			CategoryID:  3,
		}

		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Post("/product", productHandler.CreateProduct)

		productService.On("CreateProduct",product).Return(nil)
		
		reqBody := []byte(`{"name":"CPU","description":"descriptionTest","Price":100,"image":"https://www.google.com","stock":10,"isFeatured":false,"isOnSale":false,"categoryID":3,}`)

		req := httptest.NewRequest("POST","/product", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Invalid request body")
	})

	t.Run("Validation error", func(t *testing.T) {
		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Post("/product", productHandler.CreateProduct)

		tests := []struct{
			name string
			reqBody []byte
			expect []string
		}{
			{
				name: "Name is required",
				reqBody: []byte(`{"name":"","description":"descriptionTest","Price":100,"image":"https://www.google.com","stock":10,"isFeatured":false,"isOnSale":false,"categoryID":3}`),
				expect: []string{"Name is required"},
			},
			{
				name: "Name is min",
				reqBody: []byte(`{"name":"t","description":"descriptionTest","Price":100,"image":"https://www.google.com","stock":10,"isFeatured":false,"isOnSale":false,"categoryID":3}`),
				expect: []string{"Name is min"},
			},
			{
				name: "Description is min",
				reqBody: []byte(`{"name":"test1","description":"test","Price":100,"image":"https://www.google.com","stock":10,"isFeatured":false,"isOnSale":false,"categoryID":3}`),
				expect: []string{"Description is min"},
			},
			{
				name: "Description is required",
				reqBody: []byte(`{"name":"test1","description":"","Price":100,"image":"https://www.google.com","stock":10,"isFeatured":false,"isOnSale":false,"categoryID":3}`),
				expect: []string{"Description is required"},
			},
		}

		for _,tc := range tests {
			t.Run(tc.name,func(t *testing.T) {
				req := httptest.NewRequest("POST","/product", bytes.NewReader(tc.reqBody))

				req.Header.Set("Content-Type", "application/json")

				res,err := app.Test(req)

				assert.NoError(t, err)
				assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

				body, _ := io.ReadAll(res.Body)

				for _, msg := range tc.expect {
					assert.Contains(t,string(body),msg)
				}
			})
		}
		
		

		

		
	})

	t.Run("Error to CreateProduct", func(t *testing.T) {
		product := &models.Product{
			Name:        "CPU",
			Description: "descriptionTest",
			Price:       100,
			Image:       "https://www.google.com",
			Stock:       10,
			IsFeatured:  false,
			IsOnSale:    false,
			CategoryID:  3,
		}

		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Post("/product", productHandler.CreateProduct)

		productService.On("CreateProduct",product).Return(errors.New("can't to create product"))
		
		reqBody := []byte(`{"name":"CPU","description":"descriptionTest","Price":100,"image":"https://www.google.com","stock":10,"isFeatured":false,"isOnSale":false,"categoryID":3}`)

		req := httptest.NewRequest("POST","/product", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"can't to create product")
	})
}

func TestUpdateProduct(t *testing.T) {
	t.Run("UpdateProduct Success", func(t *testing.T) {
		product := &models.Product{
			Name:        "CPU Update",
			Description: "descriptionTest",
			Price:       100,
			Image:       "https://www.google.com",
			Stock:       10,
			IsFeatured:  false,
			IsOnSale:    false,
			CategoryID:  3,
		}

		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Put("/product/:id", productHandler.UpdateProduct)

		productService.On("UpdateProduct",uint(1),product).Return(nil)
		
		reqBody := []byte(`{"name":"CPU Update","description":"descriptionTest","Price":100,"image":"https://www.google.com","stock":10,"isFeatured":false,"isOnSale":false,"categoryID":3}`)

		req := httptest.NewRequest("PUT","/product/1", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Product updated successfully")
	})

	t.Run("Invalid product ID", func(t *testing.T) {
	app := fiber.New()

	productService := services.NewProductServiceMock()
	productHandler := handlers.NewProductHandler(productService)

	app.Put("/product/:id", productHandler.UpdateProduct)

	req := httptest.NewRequest("PUT", "/product/abc", bytes.NewReader(nil))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	assert.Contains(t, string(body), "Invalid product ID")

	productService.AssertExpectations(t)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		product := &models.Product{
			Name:        "CPU Update",
			Description: "descriptionTest",
			Price:       100,
			Image:       "https://www.google.com",
			Stock:       10,
			IsFeatured:  false,
			IsOnSale:    false,
			CategoryID:  3,
		}

		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Put("/product/:id", productHandler.UpdateProduct)

		productService.On("UpdateProduct",uint(1),product).Return(nil)
		
		reqBody := []byte(`{"name":"CPU Update","description":"descriptionTest","Price":100,"image":"https://www.google.com","stock":10,"isFeatured":false,"isOnSale":false,"categoryID":3,}`) // NOTE - ใส่ , เกิน
 
		req := httptest.NewRequest("PUT","/product/1", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Invalid request body")
	})

	t.Run("Validation error",func(t *testing.T) {
		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Put("/product/:id", productHandler.UpdateProduct)

		tests := []struct{
			name string
			reqBody []byte
			expect []string
		}{
			{
				name: "Name is required",
				reqBody: []byte(`{"name":"","description":"descriptionTest","Price":100,"image":"https://www.google.com","stock":10,"isFeatured":false,"isOnSale":false,"categoryID":3}`),
				expect: []string{"Name is required"},
			},
			{
				name: "Name is min",
				reqBody: []byte(`{"name":"t","description":"descriptionTest","Price":100,"image":"https://www.google.com","stock":10,"isFeatured":false,"isOnSale":false,"categoryID":3}`),
				expect: []string{"Name is min"},
			},
			{
				name: "Description is min",
				reqBody: []byte(`{"name":"test1","description":"test","Price":100,"image":"https://www.google.com","stock":10,"isFeatured":false,"isOnSale":false,"categoryID":3}`),
				expect: []string{"Description is min"},
			},
			{
				name: "Description is required",
				reqBody: []byte(`{"name":"test1","description":"","Price":100,"image":"https://www.google.com","stock":10,"isFeatured":false,"isOnSale":false,"categoryID":3}`),
				expect: []string{"Description is required"},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name,func(t *testing.T) {
				req := httptest.NewRequest("PUT","/product/1", bytes.NewReader(tc.reqBody))
				req.Header.Set("Content-Type", "application/json")

				res,err := app.Test(req)
				assert.NoError(t, err)
				assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

				body, _ := io.ReadAll(res.Body)

				for _, msg := range tc.expect {
					assert.Contains(t,string(body),msg)
				}
			})
		}

		

	

		
	})

	t.Run("Error to updateProduct", func(t *testing.T) {
		product := &models.Product{
			Name:        "CPU Update",
			Description: "descriptionTest",
			Price:       100,
			Image:       "https://www.google.com",
			Stock:       10,
			IsFeatured:  false,
			IsOnSale:    false,
			CategoryID:  3,
		}

		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Put("/product/:id", productHandler.UpdateProduct)

		productService.On("UpdateProduct",uint(1),product).Return(errors.New("Fail to update Product"))
		
		reqBody := []byte(`{"name":"CPU Update","description":"descriptionTest","Price":100,"image":"https://www.google.com","stock":10,"isFeatured":false,"isOnSale":false,"categoryID":3}`)

		req := httptest.NewRequest("PUT","/product/1", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Fail to update Product")
	})
}

func TestDeleteProduct(t *testing.T) {
	t.Run("DeleteProduct Success", func(t *testing.T) {
		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Delete("/product/:id", productHandler.DeleteProduct)

		productService.On("DeleteProduct",uint(1)).Return(nil)

		req := httptest.NewRequest("DELETE","/product/1", bytes.NewReader(nil))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Product deleted successfully")
	})

	t.Run("Invalid product ID", func(t *testing.T) {
		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Delete("/product/:id", productHandler.DeleteProduct)

		productService.On("DeleteProduct",uint(1)).Return(nil)

		req := httptest.NewRequest("DELETE","/product/fake", bytes.NewReader(nil))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Invalid product ID")
	})

	t.Run("Error Delete Product", func(t *testing.T) {
		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Delete("/product/:id", productHandler.DeleteProduct)

		productService.On("DeleteProduct",uint(1)).Return(errors.New("Error delete product"))

		req := httptest.NewRequest("DELETE","/product/1", bytes.NewReader(nil))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Error delete product")
	})
}

func TestGetProductById(t *testing.T) {
	t.Run("GetProductByID Success",func(t *testing.T) {
		product := &models.Product{
			Name:        "CPU",
			Description: "descriptionTest",
			Price:       100,
			Image:       "https://www.google.com",
			Stock:       10,
			IsFeatured:  false,
			IsOnSale:    false,
			CategoryID:  3,
		}

		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Get("/product/:id", productHandler.GetProductByID)

		productService.On("GetProductByID",uint(1)).Return(product,nil)

		req := httptest.NewRequest("GET","/product/1", bytes.NewReader(nil))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Product retrieved successfully")
	})

	t.Run("Invalid product ID",func(t *testing.T) {
		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Get("/product/:id", productHandler.GetProductByID)

		req := httptest.NewRequest("GET","/product/fake", bytes.NewReader(nil))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Invalid product ID")
	})

	t.Run("GetProductByID from service error",func(t *testing.T) {

		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Get("/product/:id", productHandler.GetProductByID)

		productService.On("GetProductByID",uint(1)).Return(nil,errors.New("error to get product"))

		req := httptest.NewRequest("GET","/product/1", bytes.NewReader(nil))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"error to get product")
	})

	t.Run("GetProductByID service not get product",func(t *testing.T) {
		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Get("/product/:id", productHandler.GetProductByID)

		productService.On("GetProductByID",uint(1)).Return(nil,nil)

		req := httptest.NewRequest("GET","/product/1", bytes.NewReader(nil))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Product not found")
	})
}

func TestGetAllProducts(t *testing.T) {
	t.Run("GetAllProducts Success",func(t *testing.T) {
		product := &models.Product{
			Name:        "CPU",
			Description: "descriptionTest",
			Price:       100,
			Image:       "https://www.google.com",
			Stock:       10,
			IsFeatured:  false,
			IsOnSale:    false,
			CategoryID:  3,
		}

		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Get("/product/", productHandler.GetAllProducts)

		minPrice := 0
		maxPrice := 999999
		searchName := ""
		category := ""
		productService.On("GetAllProducts",uint(1),uint(10),int64(minPrice),int64(maxPrice),searchName,category).Return(product,int64(1),nil)

		req := httptest.NewRequest("GET","/product", bytes.NewReader(nil))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Products retrieved successfully")
	})

	t.Run("Limit less than 1",func(t *testing.T) {
		product := &models.Product{
			Name:        "CPU",
			Description: "descriptionTest",
			Price:       100,
			Image:       "https://www.google.com",
			Stock:       10,
			IsFeatured:  false,
			IsOnSale:    false,
			CategoryID:  3,
		}

		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Get("/product/", productHandler.GetAllProducts)

		minPrice := 0
		maxPrice := 999999
		searchName := ""
		category := ""
		productService.On("GetAllProducts",uint(1),uint(10),int64(minPrice),int64(maxPrice),searchName,category).Return(product,int64(1),nil)

		req := httptest.NewRequest("GET","/product?page=1&limit=0", bytes.NewReader(nil))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Products retrieved successfully")
	})

	t.Run("MinPrice more than MaxPrice",func(t *testing.T) {
		product := &models.Product{
			Name:        "CPU",
			Description: "descriptionTest",
			Price:       100,
			Image:       "https://www.google.com",
			Stock:       10,
			IsFeatured:  false,
			IsOnSale:    false,
			CategoryID:  3,
		}

		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Get("/product/", productHandler.GetAllProducts)

		minPrice := 200
		maxPrice := 100
		searchName := ""
		category := ""
		productService.On("GetAllProducts",uint(1),uint(10),int64(minPrice),int64(maxPrice),searchName,category).Return(product,int64(1),nil)

		req := httptest.NewRequest("GET","/product?page=1&limit=10&minPrice=200&maxPrice=100", bytes.NewReader(nil))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"minPrice must be less than maxPrice")
	})

	t.Run("Error GetAllProducts",func(t *testing.T) {
		app := fiber.New()

		productService := services.NewProductServiceMock()
		productHandler := handlers.NewProductHandler(productService)

		app.Get("/product/", productHandler.GetAllProducts)

		minPrice := 0
		maxPrice := 999999
		searchName := ""
		category := ""
		productService.On("GetAllProducts",uint(1),uint(10),int64(minPrice),int64(maxPrice),searchName,category).Return(nil,int64(0),errors.New("Not get product"))

		req := httptest.NewRequest("GET","/product", bytes.NewReader(nil))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Not get product")
	})
}