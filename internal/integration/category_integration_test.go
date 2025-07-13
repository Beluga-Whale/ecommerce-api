package integration_test

import (
	"bytes"
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

func setUpAppCategory() *fiber.App {
	// NOTE - LoadEnv
	config.LoadEnv()

	// NOTE - Connect DB
	config.ConnectTestDB()
	hashPassword := utils.NewPasswordUtil()
	jwtUtil := utils.NewJwt()

	categoryRepo := repositories.NewCategoryRepository(config.TestDB)
	categoryService := services.NewCategoryService(categoryRepo)

	userRepo := repositories.NewUserRepository(config.TestDB)
	userService := services.NewUserService(userRepo,hashPassword,jwtUtil)
	userHandler := handlers.NewUserHandler(userService)

	categoryHandler := handlers.NewCategoryHandler(categoryService)
	// NOTE - Fiber
	app := fiber.New()

	app.Post("/register", userHandler.Register)
	app.Post("/login",userHandler.Login)
	app.Get("/category", categoryHandler.GetAll)
	app.Post("/category",middleware.AuthMiddleware(jwtUtil), categoryHandler.Create)
	app.Put("/category/:id",middleware.AuthMiddleware(jwtUtil), categoryHandler.Update)
	app.Delete("/category/:id",middleware.AuthMiddleware(jwtUtil), categoryHandler.Delete)
	
	return app
}

func clearDataBaseUserCategory(){
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

func RegisterAndLoginCategory(t *testing.T,app *fiber.App, email string, password string) string {
		
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

func TestCrateCategoryIntegration(t *testing.T) {
	t.Run("Integration Create Success", func(t *testing.T) {
		app := setUpAppCategory()

		email := "halay@gmail.com"
		password := "password"

		token := RegisterAndLoginCategory(t, app, email, password)

		reqBody := []byte(`{"name":"Sweater"}`)

		req := httptest.NewRequest("POST", "/category", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		assert.Contains(t, string(body), "Category created successfully")

		clearDataBaseUserCategory()
	})

	t.Run("Integration Invalid request body", func(t *testing.T) {
		app := setUpAppCategory()

		email := "halay@gmail.com"
		password := "password"

		token := RegisterAndLoginCategory(t, app, email, password)

		reqBody := []byte(`{InvalidBody}`)

		req := httptest.NewRequest("POST", "/category", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		res, err := app.Test(req)
		
		body, _ := io.ReadAll(res.Body)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		assert.Contains(t, string(body), "Invalid request body")
		clearDataBaseUserCategory()
	})

	t.Run("Integration RequestBody Is Empty", func(t *testing.T) {
		app := setUpAppCategory()

		email := "halay@gmail.com"
		password := "password"

		token := RegisterAndLoginCategory(t, app, email, password)

		reqBody := []byte(`{"name":""}`)

		req := httptest.NewRequest("POST", "/category", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		res, err := app.Test(req)
		
		body, _ := io.ReadAll(res.Body)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		assert.Contains(t, string(body), "Name is required")
		clearDataBaseUserCategory()
	})

	t.Run("Integration Already Exist Category", func(t *testing.T) {
		app := setUpAppCategory()

		email := "halay@gmail.com"
		password := "password"

		token := RegisterAndLoginCategory(t, app, email, password)

		reqBody := []byte(`{"name":"Sweater"}`)

		req := httptest.NewRequest("POST", "/category", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		// NOTE สร้างอีกตัวนึงที่ชื่อเหมือนกัน
		reqBody = []byte(`{"name":"Sweater"}`)

		req = httptest.NewRequest("POST", "/category", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ = io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		assert.Contains(t, string(body), "Category already exists")

		clearDataBaseUserCategory()
	})
	
}

func TestUpdateCategoryIntegration(t *testing.T) {
	t.Run("Integration Update Success",func(t *testing.T) {
		app := setUpAppCategory()

		email := "halay@gmail.com"
		password := "password"

		token := RegisterAndLoginCategory(t, app, email, password)

		reqBody := []byte(`{"name":"Sweater"}`)

		req := httptest.NewRequest("POST", "/category", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode, "Expected 201 Created")

		assert.Contains(t, string(body), "Category created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var categoryId models.Category
		if err := config.TestDB.Where("name = ?","Sweater").First(&categoryId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}

		// NOTE - UPDATE 
		
		reqBodyUpdate := []byte(`{"name":"SweaterUpdate"}`)
		req = httptest.NewRequest("PUT",  fmt.Sprintf("/category/%d", categoryId.ID), bytes.NewReader(reqBodyUpdate))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusOK, res.StatusCode, "Category update successfully")

		clearDataBaseUserCategory()
	})

	t.Run("Integration Update InvalidRequestBody",func(t *testing.T) {
		app := setUpAppCategory()

		email := "halay@gmail.com"
		password := "password"

		token := RegisterAndLoginCategory(t, app, email, password)

		reqBody := []byte(`{"name":"Sweater"}`)

		req := httptest.NewRequest("POST", "/category", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode, "Expected 201 Created")

		assert.Contains(t, string(body), "Category created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var categoryId models.Category
		if err := config.TestDB.Where("name = ?","Sweater").First(&categoryId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}

		// NOTE - UPDATE 
		
		reqBodyUpdate := []byte(`Invalid Body`)
		req = httptest.NewRequest("PUT",  fmt.Sprintf("/category/%d", categoryId.ID), bytes.NewReader(reqBodyUpdate))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode, "Invalid request body")

		clearDataBaseUserCategory()
	})

	t.Run("Integration Update Name Is Empty",func(t *testing.T) {
		app := setUpAppCategory()

		email := "halay@gmail.com"
		password := "password"

		token := RegisterAndLoginCategory(t, app, email, password)

		reqBody := []byte(`{"name":"Sweater"}`)

		req := httptest.NewRequest("POST", "/category", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode, "Expected 201 Created")

		assert.Contains(t, string(body), "Category created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var categoryId models.Category
		if err := config.TestDB.Where("name = ?","Sweater").First(&categoryId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}

		// NOTE - UPDATE 
		
		reqBodyUpdate := []byte(`{
			"name":""
		}`)
		req = httptest.NewRequest("PUT",  fmt.Sprintf("/category/%d", categoryId.ID), bytes.NewReader(reqBodyUpdate))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode, "Name is required")
		clearDataBaseUserCategory()
	})

	t.Run("Integration Update Exist Category",func(t *testing.T) {
		app := setUpAppCategory()

		email := "halay@gmail.com"
		password := "password"

		token := RegisterAndLoginCategory(t, app, email, password)

		reqBody := []byte(`{"name":"Sweater"}`)

		req := httptest.NewRequest("POST", "/category", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode, "Expected 201 Created")

		assert.Contains(t, string(body), "Category created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var categoryId models.Category
		if err := config.TestDB.Where("name = ?","Sweater").First(&categoryId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}

		// NOTE - UPDATE 
		
		reqBodyUpdate := []byte(`{"name":"SweaterUpdate"}`)
		req = httptest.NewRequest("PUT",  fmt.Sprintf("/category/%d", uint(1000)), bytes.NewReader(reqBodyUpdate))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		assert.Contains(t,string(body), "Category not found")

		clearDataBaseUserCategory()
	})
}

func TestDeleteIntegration(t *testing.T) {
	t.Run("Integration Delete Success",func(t *testing.T) {
		app := setUpAppCategory()

		email := "halay@gmail.com"
		password := "password"

		token := RegisterAndLoginCategory(t, app, email, password)

		reqBody := []byte(`{"name":"Sweater"}`)

		req := httptest.NewRequest("POST", "/category", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode, "Expected 201 Created")

		assert.Contains(t, string(body), "Category created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var categoryId models.Category
		if err := config.TestDB.Where("name = ?","Sweater").First(&categoryId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}

		// NOTE - Delete
		
		req = httptest.NewRequest("DELETE",  fmt.Sprintf("/category/%d", categoryId.ID),nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.Contains(t, string(body), "Category deleted successfully")
		clearDataBaseUserCategory()
	})

	t.Run("Integration Delete Non Exist  Category",func(t *testing.T) {
		app := setUpAppCategory()

		email := "halay@gmail.com"
		password := "password"

		token := RegisterAndLoginCategory(t, app, email, password)

		reqBody := []byte(`{"name":"Sweater"}`)

		req := httptest.NewRequest("POST", "/category", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode, "Expected 201 Created")

		assert.Contains(t, string(body), "Category created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var categoryId models.Category
		if err := config.TestDB.Where("name = ?","Sweater").First(&categoryId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}

		// NOTE - Delete
		
		req = httptest.NewRequest("DELETE",  fmt.Sprintf("/category/%d", uint(999)),nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		assert.Contains(t, string(body), "Category not found")
		clearDataBaseUserCategory()
	})

	t.Run("Integration Delete Unauthorized ",func(t *testing.T) {
		app := setUpAppCategory()

		email := "halay@gmail.com"
		password := "password"

		token := RegisterAndLoginCategory(t, app, email, password)

		reqBody := []byte(`{"name":"Sweater"}`)

		req := httptest.NewRequest("POST", "/category", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt="+token)

		res, err := app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusCreated, res.StatusCode, "Expected 201 Created")

		assert.Contains(t, string(body), "Category created successfully")

		//NOTE SELECT id ของ category มีก่อน
		var categoryId models.Category
		if err := config.TestDB.Where("name = ?","Sweater").First(&categoryId).Error;err !=nil{
			t.Fatalf("Failed to fetch category: %v", err)
		}

		// NOTE - Delete
		
		req = httptest.NewRequest("DELETE",  fmt.Sprintf("/category/%d", categoryId.ID),nil)
		req.Header.Set("Content-Type", "application/json")
		
		res, err = app.Test(req)
		require.NoError(t, err, "Request failed")

		body, _= io.ReadAll(res.Body)

		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
		assert.Contains(t, string(body), "Unauthorized - missing token")
		clearDataBaseUserCategory()
	})
}
