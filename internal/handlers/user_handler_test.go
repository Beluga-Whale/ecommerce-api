package handlers_test

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Beluga-Whale/ecommerce-api/internal/handlers"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T){
	t.Run("Test Register Success",func(t *testing.T) {
		user := &models.User{
			Email: "test@gmail.com",
			Password: "password123",
			Name: "Test User",
		}

		userService := services.NewUserServiceMock()

		userService.On("Register", user).Return(nil)

		userHandler := handlers.NewUserHandler(userService)

		app := fiber.New()
		app.Post("/register", userHandler.Register)

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
	} )
}