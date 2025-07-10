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
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestRegister(t *testing.T) {
	t.Run("Register Success",func(t *testing.T) {
		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		userService.On("Register",mock.Anything).Return(nil)

		app := fiber.New()
		app.Post("/register",userHandler.Register)
		
		reqBody:= []byte(`{
			"email":"test@gmail.com",
			"password":"password",
			"firstName":"FirstName",
			"lastName":"LastName",
			"phone":"0938654678",
			"birthDate":"2011-10-05T00:00:00.000Z"
		}`)

		req :=httptest.NewRequest("POST","/register",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "User registered success")
	})

	t.Run("Invalid Request body",func(t *testing.T) {
		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		userService.On("Register",mock.Anything).Return(nil)

		app := fiber.New()
		app.Post("/register",userHandler.Register)
		
		req :=httptest.NewRequest("POST","/register",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request body")
	})

	t.Run("Validate Request Body",func(t *testing.T) {
		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		userService.On("Register",mock.Anything).Return(nil)

		app := fiber.New()
		app.Post("/register",userHandler.Register)
		
		reqBody:= []byte(`{
			"password":"password",
			"firstName":"FirstName",
			"lastName":"LastName",
			"phone":"0938654678",
			"birthDate":"2011-10-05T00:00:00.000Z"
		}`)

		req :=httptest.NewRequest("POST","/register",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Email is required")
	})

	t.Run("Error To Register",func(t *testing.T) {
		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		userService.On("Register",mock.Anything).Return(errors.New("User registered successfully"))

		app := fiber.New()
		app.Post("/register",userHandler.Register)
		
		reqBody:= []byte(`{
			"email":"test@gmail.com",
			"password":"password",
			"firstName":"FirstName",
			"lastName":"LastName",
			"phone":"0938654678",
			"birthDate":"2011-10-05T00:00:00.000Z"
		}`)

		req :=httptest.NewRequest("POST","/register",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "User registered successfully")
	})

}

func TestLogin(t *testing.T) {
	t.Run("Login Success",func(t *testing.T) {
		jwtToken := "fake_jwtToken"

		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		userService.On("Login",mock.Anything).Return(jwtToken,1,nil)

		app := fiber.New()
		app.Post("/login",userHandler.Login)
		
		reqBody:= []byte(`{
			"email":"test@gmail.com",
			"password":"password"
		}`)

		req :=httptest.NewRequest("POST","/login",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Login successful")
	})
	t.Run("Invalid request body",func(t *testing.T) {
		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		app := fiber.New()
		app.Post("/login",userHandler.Login)
		
		req :=httptest.NewRequest("POST","/login",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request body")
	})

	t.Run("Validate request body",func(t *testing.T) {

		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		app := fiber.New()
		app.Post("/login",userHandler.Login)
		
		reqBody:= []byte(`{
			"password":"password"
		}`)

		req :=httptest.NewRequest("POST","/login",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Email is required")
	})
	t.Run("Error to login",func(t *testing.T) {
		jwtToken := "fake_jwtToken"

		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		userService.On("Login",mock.Anything).Return(jwtToken,1,errors.New("Error to login"))

		app := fiber.New()
		app.Post("/login",userHandler.Login)
		
		reqBody:= []byte(`{
			"email":"test@gmail.com",
			"password":"password"
		}`)

		req :=httptest.NewRequest("POST","/login",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Error to login")
	})
}

func TestGetProfile(t *testing.T) {
	t.Run("GeProfile Success",func(t *testing.T) {
		userMock := models.User{
			Model: gorm.Model{ID: 1},
			Email: "test@gmail.com",
		}

		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		userService.On("GetProfile",mock.Anything).Return(&userMock,nil)

		app := fiber.New()
		app.Get("/user/profile",testMiddleware,userHandler.GetProfile)

		req :=httptest.NewRequest("GET","/user/profile",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "test@gmail.com")
	})

	t.Run("Not have UserID",func(t *testing.T) {
		userMock := models.User{
			Model: gorm.Model{ID: 1},
			Email: "test@gmail.com",
		}

		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", nil)
			return c.Next()
		}

		userService.On("GetProfile",mock.Anything).Return(&userMock,nil)

		app := fiber.New()
		app.Get("/user/profile",testMiddleware,userHandler.GetProfile)

		req :=httptest.NewRequest("GET","/user/profile",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Unauthorized")
	})

	t.Run("Invalid user ID forat",func(t *testing.T) {
		userMock := models.User{
			Model: gorm.Model{ID: 1},
			Email: "test@gmail.com",
		}

		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "testFail")
			return c.Next()
		}

		userService.On("GetProfile",mock.Anything).Return(&userMock,nil)

		app := fiber.New()
		app.Get("/user/profile",testMiddleware,userHandler.GetProfile)

		req :=httptest.NewRequest("GET","/user/profile",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid user ID format")
	})
}

func TestUpdateProfile(t *testing.T) {
	t.Run("Update Profile Success",func(t *testing.T) {

		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		userService.On("UpdateProfile",mock.Anything,mock.Anything).Return(nil)

		app := fiber.New()
		app.Patch("/user/profile",testMiddleware,userHandler.UpdateProfile)

		reqBody:= []byte(`{
			"firstName":"UpdateFirstName"
		}`)


		req :=httptest.NewRequest("PATCH","/user/profile",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Profile updated successfully")
	})

	t.Run("Invalid request body",func(t *testing.T) {

		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		userService.On("UpdateProfile",mock.Anything,mock.Anything).Return(nil)

		app := fiber.New()
		app.Patch("/user/profile",testMiddleware,userHandler.UpdateProfile)

		req :=httptest.NewRequest("PATCH","/user/profile",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request body")
	})

	t.Run("Not have UserID",func(t *testing.T) {

		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", nil)
			return c.Next()
		}

		userService.On("UpdateProfile",mock.Anything,mock.Anything).Return(nil)

		app := fiber.New()
		app.Patch("/user/profile",testMiddleware,userHandler.UpdateProfile)

		reqBody:= []byte(`{
			"firstName":"UpdateFirstName"
		}`)


		req :=httptest.NewRequest("PATCH","/user/profile",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Unauthorized")
	})

	t.Run("Update Profile Success",func(t *testing.T) {

		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "sdfsf")
			return c.Next()
		}

		userService.On("UpdateProfile",mock.Anything,mock.Anything).Return(nil)

		app := fiber.New()
		app.Patch("/user/profile",testMiddleware,userHandler.UpdateProfile)

		reqBody:= []byte(`{
			"firstName":"UpdateFirstName"
		}`)


		req :=httptest.NewRequest("PATCH","/user/profile",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid user ID format")
	})

	t.Run("Update Profile Success",func(t *testing.T) {

		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		userService.On("UpdateProfile",mock.Anything,mock.Anything).Return(errors.New("Failed to update profile"))

		app := fiber.New()
		app.Patch("/user/profile",testMiddleware,userHandler.UpdateProfile)

		reqBody:= []byte(`{
			"firstName":"UpdateFirstName"
		}`)


		req :=httptest.NewRequest("PATCH","/user/profile",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Failed to update profile")
	})
}

func TestLogout(t *testing.T) {
	t.Run("Logout",func(t *testing.T) {
		userService := services.NewUserServiceMock()

		userHandler := handlers.NewUserHandler(userService)

		app := fiber.New()
		app.Post("/logout",userHandler.Logout)

		req :=httptest.NewRequest("POST","/logout",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Logout successful")
	})
}