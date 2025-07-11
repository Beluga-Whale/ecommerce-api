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
	"github.com/Beluga-Whale/ecommerce-api/internal/repositories"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/Beluga-Whale/ecommerce-api/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setUpAppUser() *fiber.App {
	// NOTE - LoadEnv
	config.LoadEnv()

	// NOTE - Connect DB
	config.ConnectTestDB()

	// NOTE - Utilities
	hashPassword := utils.NewPasswordUtil()
	jwtUtil := utils.NewJwt()

	userRepo := repositories.NewUserRepository(config.TestDB)
	userService := services.NewUserService(userRepo,hashPassword,jwtUtil)

	userHandler := handlers.NewUserHandler(userService)
	// NOTE - Fiber
	app := fiber.New()

	app.Post("/register", userHandler.Register)
	app.Post("/login",userHandler.Login)
	app.Post("/logout",userHandler.Logout)
	app.Get("/user/profile",middleware.AuthMiddleware(jwtUtil),middleware.RequireRole("user"),userHandler.GetProfile)
	
	return app
}

func clearDataBaseUser(){
	if err := config.TestDB.Exec("DELETE FROM users").Error; err != nil {
		log.Fatalf("Failed to clear test database: %v", err)
	}
}


func LoginAndGetTokenUser(t *testing.T, app *fiber.App, email string, password string) string {
	reqBody := []byte(fmt.Sprintf(`{
		"email": "%s",
		"password": "%s"
	}`, email, password))

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
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

func RegisterUser(t *testing.T, email string) {
	app := setUpAppUser()
		
	reqBody := []byte(fmt.Sprintf(`{
		"email":"%s",
		"firstName":"halay2",
		"lastName":"halay1teT",
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
}

func TestRegisterIntegration(t *testing.T){
	t.Run("Integration Register Success",func(t *testing.T) {
		app := setUpAppUser()
		
		reqBody := []byte(`{
			"email":"integration@gmail.com",
			"firstName":"halay2",
			"lastName":"halay1teT",
			"password":"password",
			"phone":"0874853567",
			"birthDate":"2011-10-05T14:48:00.000Z"
		}`)
			
			req := httptest.NewRequest("POST", "/register", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			
			res,err := app.Test(req)
			
			assert.NoError(t,err)
			assert.Equal(t,fiber.StatusCreated,res.StatusCode)
			
			body,_ := io.ReadAll(res.Body)
			
			assert.Contains(t,string(body),"User registered successfully")
			clearDataBaseUser()
	})

	t.Run("Integration Register Already Exist",func(t *testing.T) {
		app := setUpAppUser()
		
		// NOTE - REgister ครั้งแรก
		reqBody := []byte(`{
			"email":"integration@gmail.com",
			"firstName":"halay2",
			"lastName":"halay1teT",
			"password":"password",
			"phone":"0874853567",
			"birthDate":"2011-10-05T14:48:00.000Z"
			}`)
			
		req := httptest.NewRequest("POST", "/register", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		
		res,err := app.Test(req)
		
		assert.NoError(t,err)
		assert.Equal(t,fiber.StatusCreated,res.StatusCode)
		
		body,_ := io.ReadAll(res.Body)
		
		assert.Contains(t,string(body),"User registered successfully")

		// NOTE -  Register คนที่ 2 ต้องไม่ได้

		reqBodyAlready := []byte(`{
			"email":"integration@gmail.com",
			"firstName":"halay2",
			"lastName":"halay1teT",
			"password":"password",
			"phone":"0874853567",
			"birthDate":"2011-10-05T14:48:00.000Z"
			}`)
			
		req = httptest.NewRequest("POST", "/register", bytes.NewReader(reqBodyAlready))
		req.Header.Set("Content-Type", "application/json")
		
		res,err = app.Test(req)
		
		assert.NoError(t,err)
		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)
		
		body,_ = io.ReadAll(res.Body)
		
		assert.Contains(t,string(body),"Email already exists")

		clearDataBaseUser()
	})

	t.Run("Integration Register InvalidBody",func(t *testing.T) {
		app := setUpAppUser()
			
		req := httptest.NewRequest("POST", "/register",nil)
		req.Header.Set("Content-Type", "application/json")
		
		res,err := app.Test(req)
		
		assert.NoError(t,err)
		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)
		
		body,_ := io.ReadAll(res.Body)
		
		assert.Contains(t,string(body),"Invalid request body")
		clearDataBaseUser()
	})

	t.Run("Integration Fail Password less than 6",func(t *testing.T) {
		app := setUpAppUser()
		
		reqBody := []byte(`{
			"email":"integration@gmail.com",
			"firstName":"halay2",
			"lastName":"halay1teT",
			"password":"pass",
			"phone":"0874853567",
			"birthDate":"2011-10-05T14:48:00.000Z"
			}`)
			
			req := httptest.NewRequest("POST", "/register", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			
			res,err := app.Test(req)
			
			assert.NoError(t,err)
			assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)
			
			body,_ := io.ReadAll(res.Body)
			
			assert.Contains(t,string(body),"Password is min")
			clearDataBaseUser()
	})
}

func TestLoginIntegration(t *testing.T){
	t.Run("Integration LoginSuccess",func(t *testing.T) {
		app := setUpAppUser()
		
		reqBody := []byte(`{
			"email":"integration@gmail.com",
			"firstName":"halay2",
			"lastName":"halay1teT",
			"password":"password",
			"phone":"0874853567",
			"birthDate":"2011-10-05T14:48:00.000Z"
			}`)
			
		req := httptest.NewRequest("POST", "/register", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		
		res,err := app.Test(req)
		
		assert.NoError(t,err)
		assert.Equal(t,fiber.StatusCreated,res.StatusCode)
		
		body,_ := io.ReadAll(res.Body)
		
		assert.Contains(t,string(body),"User registered successfully")

		reqBodyLogin := []byte(`{
			"email":"integration@gmail.com",
			"password":"password"
		}`)

		req = httptest.NewRequest("POST","/login", bytes.NewReader(reqBodyLogin))

		req.Header.Set("Content-Type", "application/json")
		// req.Header.Set("Cookie","jwt=fake-jwt-token")
		
		res,err = app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusOK,res.StatusCode)
		body,_ = io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Login successful")
		clearDataBaseUser()
	})

	t.Run("Integration Login Invalid",func(t *testing.T) {
		app := setUpAppUser()
		
		reqBody := []byte(`{
			"email":"integration@gmail.com",
			"firstName":"halay2",
			"lastName":"halay1teT",
			"password":"password",
			"phone":"0874853567",
			"birthDate":"2011-10-05T14:48:00.000Z"
			}`)
			
		req := httptest.NewRequest("POST", "/register", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		
		res,err := app.Test(req)
		
		assert.NoError(t,err)
		assert.Equal(t,fiber.StatusCreated,res.StatusCode)
		
		body,_ := io.ReadAll(res.Body)
		
		assert.Contains(t,string(body),"User registered successfully")

		req = httptest.NewRequest("POST","/login",nil)

		req.Header.Set("Content-Type", "application/json")
		// req.Header.Set("Cookie","jwt=fake-jwt-token")
		
		res,err = app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusBadRequest,res.StatusCode)
		body,_ = io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Invalid request body")
		clearDataBaseUser()
	})

	
	
}

func TestGetProfileIntegration(t *testing.T) {
	t.Run("Integration GetProfile",func(t *testing.T) {
		app := setUpAppUser()

		email := "halay@gmail.com"
		password := "password"

		RegisterUser(t, email)
		token := LoginAndGetTokenUser(t, app, email, password)

		req := httptest.NewRequest("GET", "/user/profile", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "jwt=" + token)

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Get Profile Successful")

		clearDataBaseUser()
	})

	t.Run("Integration Unauthorized To GetProfile",func(t *testing.T) {
		app := setUpAppUser()

		email := "halay@gmail.com"

		RegisterUser(t, email)

		req := httptest.NewRequest("GET", "/user/profile", nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Unauthorized")

		clearDataBaseUser()		
	})
}