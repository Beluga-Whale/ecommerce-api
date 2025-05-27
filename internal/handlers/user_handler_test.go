package handlers_test

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Beluga-Whale/ecommerce-api/internal/handlers"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
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