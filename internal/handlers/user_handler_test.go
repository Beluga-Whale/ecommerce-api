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

func TestRegister(t *testing.T){


	t.Run("Test Register Success",func(t *testing.T) {			
		app := fiber.New()
		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)
		app.Post("/register", userHandler.Register)
		userService.On("Register", mock.AnythingOfType("*models.User")).Return(nil)

		reqBody := []byte(`{
			"email":"test@gmail.com",
			"password":"password123",
			"name":"Test User"	
		}`)

		req := httptest.NewRequest("POST","/register", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"User registered successfully")
		userService.AssertExpectations(t)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		app := fiber.New()
		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)
		app.Post("/register", userHandler.Register)
		userService.On("Register", mock.AnythingOfType("*models.User")).Return(nil)

		req := httptest.NewRequest("POST", "/register", nil)
		req.Header.Set("Content-Type", "application/json")
		
		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request body")
	})

	t.Run("Validation Error", func(t *testing.T) {
		app := fiber.New()
		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)
		app.Post("/register", userHandler.Register)
		userService.On("Register", mock.AnythingOfType("*models.User")).Return(nil)

		reqBody := []byte(`{
			"email":"test",
			"password":"pa",
			"name":""	
		}`)

		req := httptest.NewRequest("POST","/register", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Email is email")
		assert.Contains(t,string(body),"Name is required")
		assert.Contains(t,string(body),"Password is min")
	})

	t.Run("Service Error", func(t *testing.T) {
		app := fiber.New()
		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)
		app.Post("/register", userHandler.Register)

		userService.On("Register", mock.AnythingOfType("*models.User")).Return(errors.New("User already exists"))

		reqBody := []byte(`{
			"email":"test@gmail.com",
			"password":"password123",
			"name":"Test User"	
		}`)

		req := httptest.NewRequest("POST","/register", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"User already exists")
		userService.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	t.Run("Test Login Success", func(t *testing.T) {
		app := fiber.New()
		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)
		app.Post("/login", userHandler.Login)

		userService.On("Login", mock.AnythingOfType("*models.User")).Return("jwt_token", nil)

		reqBody := []byte(`{
			"email":"test@gmail.com",
			"password":"password123"
		}`)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")	

		res,err := app.Test(req)	

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)

		assert.Contains(t, string(body), "jwt_token")
		assert.Contains(t, string(body), "Login successful")
	})

	t.Run("Invalid request body", func(t *testing.T) {
		app := fiber.New()
		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)
		app.Post("/login", userHandler.Login)

		reqBody := []byte(`{
			"email":"test@gmail.com",
			"password":"password123",
		}`)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")	

		res,err := app.Test(req)	

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)

		assert.Contains(t, string(body), "Invalid request body")
	})

	t.Run("Validation Error", func(t *testing.T) {
		app := fiber.New()
		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)
		app.Post("/login", userHandler.Login)

		tests := []struct{
			name string
			requestBody string
			expect []string
		}{
			{
				name :"Miss both fields",
				requestBody: `{
					"email":"",
					"password":""
				}`,
				expect: []string{
					"Email is required",
					"Password is required",
				},

			},
			{
				name :"Email is email",
				requestBody: `{
					"email": "test@gmail",
					"password": "password"
				}`,
				expect: []string{
					"Email is email",
				},
			},
			{
				name :"Password too short",
				requestBody: `{
					"email": "test@gmail.com",
					"password": "pass"
				}`,
				expect: []string{
					"Password is min",
				},
			},

		}


		for _, tc := range tests{
			t.Run(tc.name,func(t *testing.T) {
				req := httptest.NewRequest("POST", "/login", bytes.NewReader([]byte(tc.requestBody)))

				req.Header.Set("Content-Type", "application/json")	
				res,err := app.Test(req)	

				assert.NoError(t, err)
				assert.Equal(t,fiber.StatusBadRequest, res.StatusCode)

				body, _ := io.ReadAll(res.Body)

				for _, expect := range tc.expect {
					assert.Contains(t, string(body),expect)
				}

			})
			
		}


		
	})

	t.Run("Login Error", func(t *testing.T) {
		app := fiber.New()
		userService := services.NewUserServiceMock()
		userHandler := handlers.NewUserHandler(userService)
		app.Post("/login", userHandler.Login)

		userService.On("Login", mock.AnythingOfType("*models.User")).Return("", errors.New("Invalid credentials"))

		reqBody := []byte(`{
			"email":"test@gmail.com",
			"password":"password"
		}`)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")	

		res,err := app.Test(req)	

		assert.NoError(t, err)
		assert.Equal(t,fiber.StatusUnauthorized, res.StatusCode)

		body, _ := io.ReadAll(res.Body)

		assert.Contains(t, string(body), "Invalid credentials")
	})
}