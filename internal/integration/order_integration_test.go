package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Beluga-Whale/ecommerce-api/config"
	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/handlers"
	"github.com/Beluga-Whale/ecommerce-api/internal/middleware"
	"github.com/Beluga-Whale/ecommerce-api/internal/repositories"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/Beluga-Whale/ecommerce-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setUpAppOrder() *fiber.App {
	// NOTE - LoadEnv
	config.LoadEnv()

	// NOTE - Connect DB
	config.ConnectTestDB()
	productUtil := utils.NewProductUtil()
	jwtUtil := utils.NewJwt()
	hashPassword := utils.NewPasswordUtil()
	
	userRepo := repositories.NewUserRepository(config.TestDB)
	orderRepo := repositories.NewOrderRepository(config.TestDB)
	productRepo := repositories.NewProductRepository(config.TestDB)
	categoryRepo := repositories.NewCategoryRepository(config.TestDB)
	
	userService := services.NewUserService(userRepo, hashPassword, jwtUtil)
	categoryService := services.NewCategoryService(categoryRepo)
	productService:= services.NewProductService(productRepo,categoryRepo)
	orderService := services.NewOrderService(config.TestDB,orderRepo,productUtil)
	
	userHandler := handlers.NewUserHandler(userService)
	productHandler:= handlers.NewProductHandler(productService) 
	orderHandler := handlers.NewOrderHandler(orderService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// NOTE - Fiber
	app := fiber.New()

	// NOTE - User
	app.Post("/register", userHandler.Register)
	app.Post("/login",userHandler.Login)

	// NOTE - Product
	app.Get("/product",productHandler.GetAllProducts)
	app.Get("/product/:id",productHandler.GetProductByID)
	app.Post("/product",middleware.AuthMiddleware(jwtUtil),productHandler.CreateProduct)
	app.Put("/product/:id",middleware.AuthMiddleware(jwtUtil),productHandler.UpdateProduct)
	app.Delete("/product/:id",middleware.AuthMiddleware(jwtUtil),productHandler.DeleteProduct)

	// NOTE - Category
	app.Post("/category",middleware.AuthMiddleware(jwtUtil), categoryHandler.Create)

	// NOTE - Order
	app.Get("/user/order/:id", orderHandler.GetOrderByID)
	app.Post("/user/order",middleware.AuthMiddleware(jwtUtil), orderHandler.CreateOrder)
	app.Patch("/user/order",middleware.AuthMiddleware(jwtUtil), orderHandler.UpdateStatusOrder)
	app.Get("/user/order",middleware.AuthMiddleware(jwtUtil), orderHandler.GetAllOrderByUserId)
	app.Patch("/user/order/:id/status",middleware.AuthMiddleware(jwtUtil), orderHandler.UpdateOrderStatusByUser)
	app.Delete("/admin/order/:id",middleware.AuthMiddleware(jwtUtil), orderHandler.DeleteOrder)

	return app
}

func clearDataBaseOrder(){
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

func RegisterAndLoginOrder(t *testing.T,app *fiber.App, email string, password string) string {
		
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

func CreateCategoryOrder(t *testing.T, app *fiber.App,token string,name string) uint  {
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

func CreateProductOrder(t *testing.T, app *fiber.App, token string,categoryID uint) error {
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

	return nil
} 

func TestCreateOrderIntegration(t *testing.T){
	t.Run("Integration CreateOrder Success",func(t *testing.T) {
		clearDataBaseOrder()
		app := setUpAppOrder()

		email := "orderuser@gmail.com"
		password := "password"

		token := RegisterAndLoginOrder(t, app, email, password)

		categoryID := CreateCategoryOrder(t, app, token, "Clothing")

		err := CreateProductOrder(t, app, token, categoryID)
		require.NoError(t, err)

		// NOTE: ดึง Variant ID จาก DB
		var variantIDs []uint
		err = config.TestDB.
			Table("product_variants").
			Select("id").
			Scan(&variantIDs).Error

		require.NoError(t, err)

		orderPayload := dto.CreateOrderRequestDTO{
		FullName:    "Thanathat Jivapaiboonsak",
		Phone:       "0899999999",
		Address:     "123 Main St",
		Province:    "Bangkok",
		District:    "Chatuchak",
		Subdistrict: "Lat Yao",
		Zipcode:     "10900",
		Items: []dto.CreateOrderItemDTO{
				{VariantID: variantIDs[0], Quantity: 1},
			},
		}

		// NOTE - แปลงจาก struct เป็น json
		body, err := json.Marshal(orderPayload)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/user/order", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)
		resBody, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(resBody), "Order created successfully")
	})

	t.Run("Integration CreateOrder Unauthorized",func(t *testing.T) {
		clearDataBaseOrder()
		app := setUpAppOrder()

		email := "orderuser@gmail.com"
		password := "password"

		token := RegisterAndLoginOrder(t, app, email, password)

		categoryID := CreateCategoryOrder(t, app, token, "Clothing")

		err := CreateProductOrder(t, app, token, categoryID)
		require.NoError(t, err)

		// NOTE: ดึง Variant ID จาก DB
		var variantIDs []uint
		err = config.TestDB.
			Table("product_variants").
			Select("id").
			Scan(&variantIDs).Error

		require.NoError(t, err)

		orderPayload := dto.CreateOrderRequestDTO{
		FullName:    "Thanathat Jivapaiboonsak",
		Phone:       "0899999999",
		Address:     "123 Main St",
		Province:    "Bangkok",
		District:    "Chatuchak",
		Subdistrict: "Lat Yao",
		Zipcode:     "10900",
		Items: []dto.CreateOrderItemDTO{
				{VariantID: variantIDs[0], Quantity: 1},
			},
		}

		// NOTE - แปลงจาก struct เป็น json
		body, err := json.Marshal(orderPayload)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/user/order", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		res, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
		resBody, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(resBody), "Unauthorized")
	})

	t.Run("Integration CreateOrder Item Is Empty",func(t *testing.T) {
		clearDataBaseOrder()
		app := setUpAppOrder()

		email := "orderuser@gmail.com"
		password := "password"

		token := RegisterAndLoginOrder(t, app, email, password)

		categoryID := CreateCategoryOrder(t, app, token, "Clothing")

		err := CreateProductOrder(t, app, token, categoryID)
		require.NoError(t, err)

		// NOTE: ดึง Variant ID จาก DB
		var variantIDs []uint
		err = config.TestDB.
			Table("product_variants").
			Select("id").
			Scan(&variantIDs).Error

		require.NoError(t, err)

		orderPayload := dto.CreateOrderRequestDTO{
		Phone:       "0899999999",
		Address:     "123 Main St",
		Province:    "Bangkok",
		District:    "Chatuchak",
		Subdistrict: "Lat Yao",
		Zipcode:     "10900",
		Items: []dto.CreateOrderItemDTO{},
		}

		// NOTE - แปลงจาก struct เป็น json
		body, err := json.Marshal(orderPayload)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/user/order", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
		resBody, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(resBody), "no item in order")		
	})

}

func TestUpdateStatusOrderIntegration(t *testing.T){
	t.Run("Integration UpdateOrder Success",func(t *testing.T) {
		clearDataBaseOrder()
		app := setUpAppOrder()

		email := "orderuser@gmail.com"
		password := "password"

		token := RegisterAndLoginOrder(t, app, email, password)

		categoryID := CreateCategoryOrder(t, app, token, "Clothing")

		err := CreateProductOrder(t, app, token, categoryID)
		require.NoError(t, err)

		// NOTE: ดึง Variant ID จาก DB
		var variantIDs []uint
		err = config.TestDB.
			Table("product_variants").
			Select("id").
			Scan(&variantIDs).Error

		require.NoError(t, err)

		orderPayload := dto.CreateOrderRequestDTO{
		FullName:    "Thanathat Jivapaiboonsak",
		Phone:       "0899999999",
		Address:     "123 Main St",
		Province:    "Bangkok",
		District:    "Chatuchak",
		Subdistrict: "Lat Yao",
		Zipcode:     "10900",
		Items: []dto.CreateOrderItemDTO{
				{VariantID: variantIDs[0], Quantity: 1},
			},
		}

		// NOTE - แปลงจาก struct เป็น json
		body, err := json.Marshal(orderPayload)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/user/order", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)
		resBody, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(resBody), "Order created successfully")

		// NOTE -หา orderID
		var orderID uint
		err = config.TestDB.Table("orders").Select("id").Scan(&orderID).Error
		require.NoError(t, err)

		var userID uint
		err = config.TestDB.Table("users").Where("email = ?",email).Select("id").Scan(&userID).Error
		require.NoError(t, err)

		// NOTE - แปลง Strut เป็น Json
		updateReq := dto.UpdateStatusOrderDTO{
			OrderId: orderID,
			Status:  "paid",
			UserId:  userID,
		}
		updateBody, _ := json.Marshal(updateReq)

		updateRequest := httptest.NewRequest("PATCH", "/user/order", bytes.NewReader(updateBody))
		updateRequest.Header.Set("Content-Type", "application/json")
		// NOTE -STRIPE_WEBHOOK_SECRET เอามาจาก .env.test เป็น key สำหรับการ test
		updateRequest.Header.Set("Authorization", "Bearer "+os.Getenv("STRIPE_WEBHOOK_SECRET"))
		updateRequest.Header.Set("Cookie", "jwt="+token)

		updateRes, err := app.Test(updateRequest)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, updateRes.StatusCode)

		updateResBody, _ := io.ReadAll(updateRes.Body)
		assert.Contains(t, string(updateResBody), "Update Status Order Sucess")

	})

	t.Run("Integration UpdateOrder Unauthorized",func(t *testing.T) {
		clearDataBaseOrder()
		app := setUpAppOrder()

		email := "orderuser@gmail.com"
		password := "password"

		token := RegisterAndLoginOrder(t, app, email, password)

		categoryID := CreateCategoryOrder(t, app, token, "Clothing")

		err := CreateProductOrder(t, app, token, categoryID)
		require.NoError(t, err)

		// NOTE: ดึง Variant ID จาก DB
		var variantIDs []uint
		err = config.TestDB.
			Table("product_variants").
			Select("id").
			Scan(&variantIDs).Error

		require.NoError(t, err)

		orderPayload := dto.CreateOrderRequestDTO{
		FullName:    "Thanathat Jivapaiboonsak",
		Phone:       "0899999999",
		Address:     "123 Main St",
		Province:    "Bangkok",
		District:    "Chatuchak",
		Subdistrict: "Lat Yao",
		Zipcode:     "10900",
		Items: []dto.CreateOrderItemDTO{
				{VariantID: variantIDs[0], Quantity: 1},
			},
		}

		// NOTE - แปลงจาก struct เป็น json
		body, err := json.Marshal(orderPayload)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/user/order", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)
		resBody, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(resBody), "Order created successfully")

		// NOTE -หา orderID
		var orderID uint
		err = config.TestDB.Table("orders").Select("id").Scan(&orderID).Error
		require.NoError(t, err)

		var userID uint
		err = config.TestDB.Table("users").Where("email = ?",email).Select("id").Scan(&userID).Error
		require.NoError(t, err)

		// NOTE - แปลง Strut เป็น Json
		updateReq := dto.UpdateStatusOrderDTO{
			OrderId: orderID,
			Status:  "paid",
			UserId:  userID,
		}
		updateBody, _ := json.Marshal(updateReq)

		updateRequest := httptest.NewRequest("PATCH", "/user/order", bytes.NewReader(updateBody))
		updateRequest.Header.Set("Content-Type", "application/json")
		// NOTE -STRIPE_WEBHOOK_SECRET เอามาจาก .env.test เป็น key สำหรับการ test
		updateRequest.Header.Set("Authorization", "Bearer "+os.Getenv("STRIPE_WEBHOOK_SECRET"))
		
		updateRes, err := app.Test(updateRequest)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, updateRes.StatusCode)

		updateResBody, _ := io.ReadAll(updateRes.Body)
		assert.Contains(t, string(updateResBody), "Unauthorized")

	})

	t.Run("Integration UpdateOrder Invalid request body",func(t *testing.T) {
		clearDataBaseOrder()
		app := setUpAppOrder()

		email := "orderuser@gmail.com"
		password := "password"

		token := RegisterAndLoginOrder(t, app, email, password)

		categoryID := CreateCategoryOrder(t, app, token, "Clothing")

		err := CreateProductOrder(t, app, token, categoryID)
		require.NoError(t, err)

		// NOTE: ดึง Variant ID จาก DB
		var variantIDs []uint
		err = config.TestDB.
			Table("product_variants").
			Select("id").
			Scan(&variantIDs).Error

		require.NoError(t, err)

		orderPayload := dto.CreateOrderRequestDTO{
		FullName:    "Thanathat Jivapaiboonsak",
		Phone:       "0899999999",
		Address:     "123 Main St",
		Province:    "Bangkok",
		District:    "Chatuchak",
		Subdistrict: "Lat Yao",
		Zipcode:     "10900",
		Items: []dto.CreateOrderItemDTO{
				{VariantID: variantIDs[0], Quantity: 1},
			},
		}

		// NOTE - แปลงจาก struct เป็น json
		body, err := json.Marshal(orderPayload)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/user/order", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)
		resBody, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(resBody), "Order created successfully")

		// NOTE -หา orderID
		var orderID uint
		err = config.TestDB.Table("orders").Select("id").Scan(&orderID).Error
		require.NoError(t, err)

		var userID uint
		err = config.TestDB.Table("users").Where("email = ?",email).Select("id").Scan(&userID).Error
		require.NoError(t, err)

		// NOTE - แปลง Strut เป็น Json
		invalidBody := []byte("this_is_not_json")

		updateRequest := httptest.NewRequest("PATCH", "/user/order", bytes.NewReader(invalidBody))
		updateRequest.Header.Set("Content-Type", "application/json")
		// NOTE -STRIPE_WEBHOOK_SECRET เอามาจาก .env.test เป็น key สำหรับการ test
		updateRequest.Header.Set("Authorization", "Bearer "+os.Getenv("STRIPE_WEBHOOK_SECRET"))
		updateRequest.Header.Set("Cookie", "jwt="+token)

		updateRes, err := app.Test(updateRequest)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, updateRes.StatusCode)

		updateResBody, _ := io.ReadAll(updateRes.Body)
		assert.Contains(t, string(updateResBody), "Invalid request body")

	})
}

func TestDeleteStatusOrderIntegration(t *testing.T){
	t.Run("Integration DeleteOrder Success",func(t *testing.T) {
		clearDataBaseOrder()
		app := setUpAppOrder()

		email := "orderuser@gmail.com"
		password := "password"

		token := RegisterAndLoginOrder(t, app, email, password)

		categoryID := CreateCategoryOrder(t, app, token, "Clothing")

		err := CreateProductOrder(t, app, token, categoryID)
		require.NoError(t, err)

		// NOTE: ดึง Variant ID จาก DB
		var variantIDs []uint
		err = config.TestDB.
			Table("product_variants").
			Select("id").
			Scan(&variantIDs).Error

		require.NoError(t, err)

		orderPayload := dto.CreateOrderRequestDTO{
		FullName:    "Thanathat Jivapaiboonsak",
		Phone:       "0899999999",
		Address:     "123 Main St",
		Province:    "Bangkok",
		District:    "Chatuchak",
		Subdistrict: "Lat Yao",
		Zipcode:     "10900",
		Items: []dto.CreateOrderItemDTO{
				{VariantID: variantIDs[0], Quantity: 1},
			},
		}

		// NOTE - แปลงจาก struct เป็น json
		body, err := json.Marshal(orderPayload)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/user/order", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)
		resBody, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(resBody), "Order created successfully")

		// NOTE -หา orderID
		var orderID uint
		err = config.TestDB.Table("orders").Select("id").Scan(&orderID).Error
		require.NoError(t, err)

		deleteRequest := httptest.NewRequest("DELETE",fmt.Sprintf("/admin/order/%d",orderID) , nil)
		deleteRequest.Header.Set("Content-Type", "application/json")
		// NOTE -STRIPE_WEBHOOK_SECRET เอามาจาก .env.test เป็น key สำหรับการ test
		deleteRequest.Header.Set("Cookie", "jwt="+token)

		updateRes, err := app.Test(deleteRequest)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, updateRes.StatusCode)

		updateResBody, _ := io.ReadAll(updateRes.Body)
		assert.Contains(t, string(updateResBody), "Order deleted successfully")

	})

	t.Run("Integration DeleteOrder Param Invalid",func(t *testing.T) {
		clearDataBaseOrder()
		app := setUpAppOrder()

		email := "orderuser@gmail.com"
		password := "password"

		token := RegisterAndLoginOrder(t, app, email, password)

		categoryID := CreateCategoryOrder(t, app, token, "Clothing")

		err := CreateProductOrder(t, app, token, categoryID)
		require.NoError(t, err)

		// NOTE: ดึง Variant ID จาก DB
		var variantIDs []uint
		err = config.TestDB.
			Table("product_variants").
			Select("id").
			Scan(&variantIDs).Error

		require.NoError(t, err)

		orderPayload := dto.CreateOrderRequestDTO{
		FullName:    "Thanathat Jivapaiboonsak",
		Phone:       "0899999999",
		Address:     "123 Main St",
		Province:    "Bangkok",
		District:    "Chatuchak",
		Subdistrict: "Lat Yao",
		Zipcode:     "10900",
		Items: []dto.CreateOrderItemDTO{
				{VariantID: variantIDs[0], Quantity: 1},
			},
		}

		// NOTE - แปลงจาก struct เป็น json
		body, err := json.Marshal(orderPayload)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/user/order", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)
		resBody, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(resBody), "Order created successfully")

		// NOTE -หา orderID
		var orderID uint
		err = config.TestDB.Table("orders").Select("id").Scan(&orderID).Error
		require.NoError(t, err)

		deleteRequest := httptest.NewRequest("DELETE","/admin/order/InvalidParams" , nil)
		deleteRequest.Header.Set("Content-Type", "application/json")
		// NOTE -STRIPE_WEBHOOK_SECRET เอามาจาก .env.test เป็น key สำหรับการ test
		deleteRequest.Header.Set("Cookie", "jwt="+token)

		updateRes, err := app.Test(deleteRequest)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, updateRes.StatusCode)

		updateResBody, _ := io.ReadAll(updateRes.Body)
		assert.Contains(t, string(updateResBody), "Invalid order ID")

	})
}

func TestGetOrderByIDIntegration(t *testing.T) {
	t.Run("Integration GetOrderByID Success",func(t *testing.T) {
		clearDataBaseOrder()
		app := setUpAppOrder()

		email := "orderuser@gmail.com"
		password := "password"

		token := RegisterAndLoginOrder(t, app, email, password)

		categoryID := CreateCategoryOrder(t, app, token, "Clothing")

		err := CreateProductOrder(t, app, token, categoryID)
		require.NoError(t, err)

		// NOTE: ดึง Variant ID จาก DB
		var variantIDs []uint
		err = config.TestDB.
			Table("product_variants").
			Select("id").
			Scan(&variantIDs).Error

		require.NoError(t, err)

		orderPayload := dto.CreateOrderRequestDTO{
		FullName:    "Thanathat Jivapaiboonsak",
		Phone:       "0899999999",
		Address:     "123 Main St",
		Province:    "Bangkok",
		District:    "Chatuchak",
		Subdistrict: "Lat Yao",
		Zipcode:     "10900",
		Items: []dto.CreateOrderItemDTO{
				{VariantID: variantIDs[0], Quantity: 1},
			},
		}

		// NOTE - แปลงจาก struct เป็น json
		body, err := json.Marshal(orderPayload)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/user/order", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)
		resBody, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(resBody), "Order created successfully")

		// NOTE -หา orderID
		var orderID uint
		err = config.TestDB.Table("orders").Select("id").Scan(&orderID).Error
		require.NoError(t, err)

		// NOTE หา userID
		var userID uint
		err = config.TestDB.Table("users").Where("email = ?",email).Select("id").Scan(&userID).Error
		require.NoError(t, err)

		getOrderByIDRequest := httptest.NewRequest("GET", fmt.Sprintf("/user/order/%d?userId=%d", orderID, userID), nil)

		getOrderByIDRequest.Header.Set("Content-Type", "application/json")

		updateRes, err := app.Test(getOrderByIDRequest)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, updateRes.StatusCode)

		updateResBody, _ := io.ReadAll(updateRes.Body)
		assert.Contains(t, string(updateResBody), "Get Order By ID Success")

	})

	t.Run("Integration GetOrderByID Invalid order ID",func(t *testing.T) {
		clearDataBaseOrder()
		app := setUpAppOrder()

		email := "orderuser@gmail.com"
		password := "password"

		token := RegisterAndLoginOrder(t, app, email, password)

		categoryID := CreateCategoryOrder(t, app, token, "Clothing")

		err := CreateProductOrder(t, app, token, categoryID)
		require.NoError(t, err)

		// NOTE: ดึง Variant ID จาก DB
		var variantIDs []uint
		err = config.TestDB.
			Table("product_variants").
			Select("id").
			Scan(&variantIDs).Error

		require.NoError(t, err)

		orderPayload := dto.CreateOrderRequestDTO{
		FullName:    "Thanathat Jivapaiboonsak",
		Phone:       "0899999999",
		Address:     "123 Main St",
		Province:    "Bangkok",
		District:    "Chatuchak",
		Subdistrict: "Lat Yao",
		Zipcode:     "10900",
		Items: []dto.CreateOrderItemDTO{
				{VariantID: variantIDs[0], Quantity: 1},
			},
		}

		// NOTE - แปลงจาก struct เป็น json
		body, err := json.Marshal(orderPayload)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/user/order", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)
		resBody, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(resBody), "Order created successfully")

		// NOTE -หา orderID
		var orderID uint
		err = config.TestDB.Table("orders").Select("id").Scan(&orderID).Error
		require.NoError(t, err)

		// NOTE หา userID
		var userID uint
		err = config.TestDB.Table("users").Where("email = ?",email).Select("id").Scan(&userID).Error
		require.NoError(t, err)

		getOrderByIDRequest := httptest.NewRequest("GET", fmt.Sprintf("/user/order/abcd?userId=%d", userID), nil)

		getOrderByIDRequest.Header.Set("Content-Type", "application/json")

		updateRes, err := app.Test(getOrderByIDRequest)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, updateRes.StatusCode)

		updateResBody, _ := io.ReadAll(updateRes.Body)
		assert.Contains(t, string(updateResBody), "Invalid order ID")

	})

	t.Run("Integration GetOrderByID Order Not Found",func(t *testing.T) {
		clearDataBaseOrder()
		app := setUpAppOrder()

		email := "orderuser@gmail.com"
		password := "password"

		token := RegisterAndLoginOrder(t, app, email, password)

		categoryID := CreateCategoryOrder(t, app, token, "Clothing")

		err := CreateProductOrder(t, app, token, categoryID)
		require.NoError(t, err)

		// NOTE: ดึง Variant ID จาก DB
		var variantIDs []uint
		err = config.TestDB.
			Table("product_variants").
			Select("id").
			Scan(&variantIDs).Error

		require.NoError(t, err)

		orderPayload := dto.CreateOrderRequestDTO{
		FullName:    "Thanathat Jivapaiboonsak",
		Phone:       "0899999999",
		Address:     "123 Main St",
		Province:    "Bangkok",
		District:    "Chatuchak",
		Subdistrict: "Lat Yao",
		Zipcode:     "10900",
		Items: []dto.CreateOrderItemDTO{
				{VariantID: variantIDs[0], Quantity: 1},
			},
		}

		// NOTE - แปลงจาก struct เป็น json
		body, err := json.Marshal(orderPayload)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/user/order", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)
		resBody, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(resBody), "Order created successfully")

		// NOTE -หา orderID
		var orderID uint
		err = config.TestDB.Table("orders").Select("id").Scan(&orderID).Error
		require.NoError(t, err)

		// NOTE หา userID
		var userID uint
		err = config.TestDB.Table("users").Where("email = ?",email).Select("id").Scan(&userID).Error
		require.NoError(t, err)


		// NOTE - Delete ก่อน เพื่อให้Order ที่จะหา หาไม่เจอ
		deleteRequest := httptest.NewRequest("DELETE",fmt.Sprintf("/admin/order/%d",orderID) , nil)
		deleteRequest.Header.Set("Content-Type", "application/json")
		// NOTE -STRIPE_WEBHOOK_SECRET เอามาจาก .env.test เป็น key สำหรับการ test
		deleteRequest.Header.Set("Cookie", "jwt="+token)

		deleteRes, err := app.Test(deleteRequest)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, deleteRes.StatusCode)

		deleteResBody, _ := io.ReadAll(deleteRes.Body)
		assert.Contains(t, string(deleteResBody), "Order deleted successfully")

		getOrderByIDRequest := httptest.NewRequest("GET", fmt.Sprintf("/user/order/%d?userId=%d", orderID, userID), nil)

		getOrderByIDRequest.Header.Set("Content-Type", "application/json")

		getOrderIdRes, err := app.Test(getOrderByIDRequest)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, getOrderIdRes.StatusCode)

		getOrderIdResBody, _ := io.ReadAll(getOrderIdRes.Body)
		assert.Contains(t, string(getOrderIdResBody), "Error go get order by ID")

	})
}