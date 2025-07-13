package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/Beluga-Whale/ecommerce-api/config"
	"github.com/Beluga-Whale/ecommerce-api/internal/handlers"
	"github.com/Beluga-Whale/ecommerce-api/internal/middleware"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/repositories"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/Beluga-Whale/ecommerce-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setUpAppProduct() *fiber.App {
	// NOTE - LoadEnv
	config.LoadEnv()

	// NOTE - Connect DB
	config.ConnectTestDB()
	jwtUtil := utils.NewJwt()
	hashPassword := utils.NewPasswordUtil()

	categoryRepo := repositories.NewCategoryRepository(config.TestDB)
	categoryService := services.NewCategoryService(categoryRepo)

	userRepo := repositories.NewUserRepository(config.TestDB)
	userService := services.NewUserService(userRepo, hashPassword, jwtUtil)
	userHandler := handlers.NewUserHandler(userService)

	categoryHandler := handlers.NewCategoryHandler(categoryService)

	productRepo := repositories.NewProductRepository(config.TestDB)
	productService:= services.NewProductService(productRepo,categoryRepo)
	productHandler:= handlers.NewProductHandler(productService) 
	// NOTE - Fiber
	app := fiber.New()

	// NOTE - User
	app.Post("/register", userHandler.Register)
	app.Post("/login",userHandler.Login)

	// NOTE - Category
	app.Post("/category",middleware.AuthMiddleware(jwtUtil), categoryHandler.Create)

	// NOTE - Product
	app.Get("/product",productHandler.GetAllProducts)
	app.Get("/product/:id",productHandler.GetProductByID)
	app.Post("/product",middleware.AuthMiddleware(jwtUtil),productHandler.CreateProduct)
	app.Put("/product/:id",middleware.AuthMiddleware(jwtUtil),productHandler.UpdateProduct)
	app.Delete("/product/:id",middleware.AuthMiddleware(jwtUtil),productHandler.DeleteProduct)

	return app
}

func clearDataBaseProduct(){
	tables := []string{
		"order_items",
		"orders",
		"payments",
		"reviews",
		"cart_items",
		"product_images",
		"product_variants",
		"products",
		"categories",
		"coupons",
		"users",
	}

	for _, table := range tables {
		if err := config.TestDB.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			log.Fatalf("Failed to clear %s: %v", table, err)
		}
	}
}

func RegisterAndLoginProduct(t *testing.T,app *fiber.App, email string, password string) string {
		
	reqBody := []byte(fmt.Sprintf(`{
		"email":"%s",
		"firstName":"halay1",
		"lastName":"halay1",
		"password":"password",
		"phone":"0874853567",
		"birthDate":"2011-10-05T14:48:00.000Z"
	}`, email))
		
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	res,err := app.Test(req)
	
	assert.NoError(t,err)
	assert.Equal(t,fiber.StatusCreated,res.StatusCode)
	
	body,_ := io.ReadAll(res.Body)
	
	assert.Contains(t,string(body),"User registered successfully")

	reqBodyLogin := []byte(fmt.Sprintf(`{
		"email": "%s",
		"password": "%s"
	}`, email, password))

	req = httptest.NewRequest("POST", "/login", bytes.NewReader(reqBodyLogin))
	req.Header.Set("Content-Type", "application/json")

	res, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	cookies := res.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "jwt" {
			return cookie.Value
		}
	}

	t.Fatal("JWT not found in cookies")
	return ""
}

func CreateCategoryProduct(t *testing.T, app *fiber.App,token string,name string) uint  {
	reqBody := []byte(fmt.Sprintf(`{
		"name":"%s"
	}`, name))

	req := httptest.NewRequest("POST", "/category", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "jwt="+token)

	res, err := app.Test(req)
	require.NoError(t, err, "Request failed")

	body, _ := io.ReadAll(res.Body)

	assert.Equal(t, fiber.StatusCreated, res.StatusCode)

	assert.Contains(t, string(body), "Category created successfully")

	var response struct {
	Data struct {
		ID uint `json:"id"`
		} `json:"data"`
	}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err, "Failed to unmarshal category create response")

	return response.Data.ID
}

func TestCreateProductIntegration(t *testing.T) {
	t.Run("Integration CreateProduct Success",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		assert.Contains(t, string(body), "Product created successfully")
		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})

	t.Run("Integration CreateProduct Invalid request body",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
			"categoryId": %d,
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		assert.Contains(t, string(body), "Invalid request body")
		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})

	t.Run("Integration CreateProduct Unauthorized",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

		assert.Contains(t, string(body), "Unauthorized")
		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})

	t.Run("Integration CreateProduct Required Field",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		assert.Contains(t, string(body), "Name is required")
		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})

	t.Run("Integration CreateProduct Image Must 3",func(t *testing.T) {
		clearDataBaseProduct()
		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
			"name": "T-Shirt",
			"title": "test title",
			"description": "<p>TEST</p>",
			"images": [
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060281/Beluga_Whale_gtyr3j.webp"
				}
			],
			"isFeatured": false,
			"isOnSale": true,
			"salePrice": 20,
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		assert.Contains(t, string(body), "You must upload exactly 3 product images")
		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})

	t.Run("Integration CreateProduct SalePrice Must More Than 0",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"salePrice": -10,
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		assert.Contains(t, string(body), "Sale price must be greater than 0")
		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})

	t.Run("Integration CreateProduct IsOnSale False But Sale Price More Than 0",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"isOnSale": false,
			"salePrice": 20,
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		assert.Contains(t, string(body), "You can should is in sale true")
		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})
}

func TestUpdateProduct(t *testing.T) {
	t.Run("Integration Update Product Success",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		assert.Contains(t, string(body), "Product created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var productId models.Product
		if err := config.TestDB.Where("name = ?","T-Shirt").First(&productId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}

		reqUpdateBody := []byte(fmt.Sprintf(`{
			"name": "T-Shirt Update",
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
			"categoryId": %d,
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
		}`, categoryID))
		req = httptest.NewRequest("PUT",  fmt.Sprintf("/product/%d", productId.ID), bytes.NewReader(reqUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusOK, res.StatusCode, "Product updated successfully")

		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})

	t.Run("Integration Update Product Unauthorized",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		assert.Contains(t, string(body), "Product created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var productId models.Product
		if err := config.TestDB.Where("name = ?","T-Shirt").First(&productId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}

		reqUpdateBody := []byte(fmt.Sprintf(`{
			"name": "T-Shirt Update",
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
			"categoryId": %d,
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
		}`, categoryID))
		req = httptest.NewRequest("PUT",  fmt.Sprintf("/product/%d", productId.ID), bytes.NewReader(reqUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
		
		assert.Contains(t, string(body), "Unauthorized")
		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})

	t.Run("Integration Update Product Invalid Product ID",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		assert.Contains(t, string(body), "Product created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var productId models.Product
		if err := config.TestDB.Where("name = ?","T-Shirt").First(&productId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}

		reqUpdateBody := []byte(fmt.Sprintf(`{
			"name": "T-Shirt Update",
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
			"categoryId": %d,
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
		}`, categoryID))
		req = httptest.NewRequest("PUT",  fmt.Sprintf("/product/%d", uint(999)), bytes.NewReader(reqUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode, "Invalid product ID")

		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})
	t.Run("Integration Update Product Required Field",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		assert.Contains(t, string(body), "Product created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var productId models.Product
		if err := config.TestDB.Where("name = ?","T-Shirt").First(&productId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}

		reqUpdateBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))
		req = httptest.NewRequest("PUT",  fmt.Sprintf("/product/%d", productId.ID), bytes.NewReader(reqUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode, "Name is required")

		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})

	t.Run("Integration Update Product Image Is Equal 3",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		assert.Contains(t, string(body), "Product created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var productId models.Product
		if err := config.TestDB.Where("name = ?","T-Shirt").First(&productId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}

		reqUpdateBody := []byte(fmt.Sprintf(`{
			"name": "T-Shirt Update",
			"title": "test title",
			"description": "<p>TEST</p>",
			"images": [
				{
					"url": "https://res.cloudinary.com/dnue94koc/image/upload/v1750060281/Beluga_Whale_gtyr3j.webp"
				}
			],
			"isFeatured": false,
			"isOnSale": true,
			"salePrice": 20,
			"categoryId": %d,
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
		}`, categoryID))
		req = httptest.NewRequest("PUT",  fmt.Sprintf("/product/%d", productId.ID), bytes.NewReader(reqUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode, "You must upload exactly 3 product images")

		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})

	t.Run("Integration Update SalePrice Must More Than 0",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		assert.Contains(t, string(body), "Product created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var productId models.Product
		if err := config.TestDB.Where("name = ?","T-Shirt").First(&productId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}

		reqUpdateBody := []byte(fmt.Sprintf(`{
			"name": "T-Shirt Update",
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
			"salePrice": -20,
			"categoryId": %d,
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
		}`, categoryID))
		req = httptest.NewRequest("PUT",  fmt.Sprintf("/product/%d", productId.ID), bytes.NewReader(reqUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode, "Sale price must be greater than 0")

		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})

	t.Run("Integration Update IsOnSale False But Sale Price More Than 0",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		assert.Contains(t, string(body), "Product created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var productId models.Product
		if err := config.TestDB.Where("name = ?","T-Shirt").First(&productId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}

		reqUpdateBody := []byte(fmt.Sprintf(`{
			"name": "T-Shirt Update",
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
			"isOnSale": false,
			"salePrice": 20,
			"categoryId": %d,
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
		}`, categoryID))
		req = httptest.NewRequest("PUT",  fmt.Sprintf("/product/%d", productId.ID), bytes.NewReader(reqUpdateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode, "You can should is in sale")

		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})
}

func TestDelete(t *testing.T) {
	t.Run("Integration Delete Success ",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		assert.Contains(t, string(body), "Product created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var productId models.Product
		if err := config.TestDB.Where("name = ?","T-Shirt").First(&productId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}
	
		req = httptest.NewRequest("DELETE",  fmt.Sprintf("/product/%d", productId.ID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusOK, res.StatusCode, "Product deleted successfully")

		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})

	t.Run("Integration Delete Invalid Product ID ",func(t *testing.T) {
			clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		assert.Contains(t, string(body), "Product created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var productId models.Product
		if err := config.TestDB.Where("name = ?","T-Shirt").First(&productId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}
	
		req = httptest.NewRequest("DELETE",  fmt.Sprintf("/product/%d", uint(999)), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode, "Invalid product ID")

		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})
}

func TestGetProductByID(t *testing.T) {
	t.Run("Integration GetProductByID Success",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		assert.Contains(t, string(body), "Product created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var productId models.Product
		if err := config.TestDB.Where("name = ?","T-Shirt").First(&productId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}
	
		req = httptest.NewRequest("GET",  fmt.Sprintf("/product/%d", productId.ID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusOK, res.StatusCode, "Product retrieved successfully")

		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})
	t.Run("Integration GetProductByID Invalid Product ID",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		assert.Contains(t, string(body), "Product created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var productId models.Product
		if err := config.TestDB.Where("name = ?","T-Shirt").First(&productId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}
	
		req = httptest.NewRequest("GET",  fmt.Sprintf("/product/%d", uint(999)), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode, "Invalid product ID")

		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})

	t.Run("Integration GetProductByID Product Not Found",func(t *testing.T) {
		clearDataBaseProduct()

		app := setUpAppProduct()

		email := "halay@gmail.com"
		password := "password"
		token := RegisterAndLoginProduct(t,app,email,password)

		categoryID := CreateCategoryProduct(t, app, token, "Clothing")
		reqBody := []byte(fmt.Sprintf(`{
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
			"categoryId": %d,
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
		}`, categoryID))

		
		req := httptest.NewRequest("POST", "/product", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		
		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		assert.Contains(t, string(body), "Product created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var productId models.Product
		if err := config.TestDB.Where("name = ?","T-Shirt").First(&productId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}
	
		req = httptest.NewRequest("GET",  fmt.Sprintf("/product/%d", productId.ID+uint(1)), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode, "Product not found")

		t.Cleanup(func() {
			clearDataBaseProduct()
		})
	})	
}
