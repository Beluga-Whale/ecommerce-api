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
	"gorm.io/gorm"
)

func TestGetAll(t *testing.T) {
	t.Run("GetAll Success",func(t *testing.T) {
		categoryMock := []models.Category{
			{
				Model: gorm.Model{ID: 1},
				Name: "TEST A",
			},
			{
				Model: gorm.Model{ID: 2},
				Name: "TEST B",
			},
		}

		categoryService := services.NewCategoryServiceMock()
		categoryService.On("GetAllCategories").Return(categoryMock,nil)

		categoryHandler := handlers.NewCategoryHandler(categoryService)

		app := fiber.New()
		app.Get("/category",categoryHandler.GetAll)

		req := httptest.NewRequest("GET", "/category",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Get all categories successfully")
	
	})

	t.Run("GetAll Success",func(t *testing.T) {

		categoryService := services.NewCategoryServiceMock()
		categoryService.On("GetAllCategories").Return(nil,errors.New("Failed to fetch categories"))

		categoryHandler := handlers.NewCategoryHandler(categoryService)

		app := fiber.New()
		app.Get("/category",categoryHandler.GetAll)

		req := httptest.NewRequest("GET", "/category",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Failed to fetch categories")
	
	})
}

func TestCreate(t *testing.T) {
	t.Run("Create category success",func(t *testing.T) {
		categoryMock := &models.Category{
			Name: "Test A",
		}

		categoryService := services.NewCategoryServiceMock()
		categoryService.On("CreateCategory",categoryMock).Return(nil)

		categoryHandler := handlers.NewCategoryHandler(categoryService)

		app := fiber.New()
		app.Post("/category",categoryHandler.Create)

		reqBody:= []byte(`{
			"name":"Test A"
		}`)

		req := httptest.NewRequest("POST", "/category",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Category created successfully")
	})

	t.Run("Invalid Request BOdy",func(t *testing.T) {
		categoryService := services.NewCategoryServiceMock()

		categoryHandler := handlers.NewCategoryHandler(categoryService)

		app := fiber.New()
		app.Post("/category",categoryHandler.Create)

		reqBody:= []byte(`fail invalid body`)

		req := httptest.NewRequest("POST", "/category",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request body")
	})

	t.Run("Required Input Name",func(t *testing.T) {
		categoryMock := &models.Category{
			Name: "Test A",
		}
		categoryService := services.NewCategoryServiceMock()
		categoryService.On("CreateCategory",categoryMock).Return(nil)

		categoryHandler := handlers.NewCategoryHandler(categoryService)

		app := fiber.New()
		app.Post("/category",categoryHandler.Create)

		reqBody:= []byte(`{
			"name":""
		}`)

		req := httptest.NewRequest("POST", "/category",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Name is required")
	})

	t.Run("Error to create category",func(t *testing.T) {
		categoryMock := &models.Category{
			Name: "Test A",
		}

		categoryService := services.NewCategoryServiceMock()
		categoryService.On("CreateCategory",categoryMock).Return(errors.New("Error to create category"))

		categoryHandler := handlers.NewCategoryHandler(categoryService)

		app := fiber.New()
		app.Post("/category",categoryHandler.Create)

		reqBody:= []byte(`{
			"name":"Test A"
		}`)

		req := httptest.NewRequest("POST", "/category",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Error to create category")
	})
}
func TestUpdate(t *testing.T) {
	t.Run("Update category success",func(t *testing.T) {
		categoryMock := &models.Category{
			Name: "Test Update",
		}

		categoryService := services.NewCategoryServiceMock()
		categoryService.On("UpdateCategory",uint(1),categoryMock).Return(nil)

		categoryHandler := handlers.NewCategoryHandler(categoryService)

		app := fiber.New()
		app.Put("/category/:id",categoryHandler.Update)

		reqBody:= []byte(`{
			"name":"Test Update"
		}`)

		req := httptest.NewRequest("PUT", "/category/1",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Category update successfully")
	})
	t.Run("Invalid Request Body",func(t *testing.T) {
		categoryMock := &models.Category{
			Name: "Test Update",
		}

		categoryService := services.NewCategoryServiceMock()
		categoryService.On("UpdateCategory",uint(1),categoryMock).Return(nil)

		categoryHandler := handlers.NewCategoryHandler(categoryService)

		app := fiber.New()
		app.Put("/category/:id",categoryHandler.Update)

		reqBody:= []byte(`Invalid Body`)

		req := httptest.NewRequest("PUT", "/category/1",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request body")
	})
	t.Run("Name Is Empty",func(t *testing.T) {
		categoryMock := &models.Category{
			Name: "Test Update",
		}

		categoryService := services.NewCategoryServiceMock()
		categoryService.On("UpdateCategory",uint(1),categoryMock).Return(nil)

		categoryHandler := handlers.NewCategoryHandler(categoryService)

		app := fiber.New()
		app.Put("/category/:id",categoryHandler.Update)

		reqBody:= []byte(`{
			"name":""
		}`)

		req := httptest.NewRequest("PUT", "/category/1",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Name is required")
	})

	t.Run("Error updateCategory",func(t *testing.T) {
		categoryMock := &models.Category{
			Name: "Test Update",
		}

		categoryService := services.NewCategoryServiceMock()
		categoryService.On("UpdateCategory",uint(1),categoryMock).Return(errors.New("Error to update category"))

		categoryHandler := handlers.NewCategoryHandler(categoryService)

		app := fiber.New()
		app.Put("/category/:id",categoryHandler.Update)

		reqBody:= []byte(`{
			"name":"Test Update"
		}`)

		req := httptest.NewRequest("PUT", "/category/1",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Error to update category")
	})

	// 	res, err := app.Test(req)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	// 	body, _ := io.ReadAll(res.Body)
	// 	assert.Contains(t, string(body), "Name is required")
	// })

	// t.Run("Error to create category",func(t *testing.T) {
	// 	categoryMock := &models.Category{
	// 		Name: "Test A",
	// 	}

	// 	categoryService := services.NewCategoryServiceMock()
	// 	categoryService.On("CreateCategory",categoryMock).Return(errors.New("Error to create category"))

	// 	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// 	app := fiber.New()
	// 	app.Post("/category",categoryHandler.Create)
}

func TestDelete(t *testing.T) {
	t.Run("Delete Success",func(t *testing.T) {

		categoryService := services.NewCategoryServiceMock()
		categoryService.On("DeleteCategory",uint(1)).Return(nil)

		categoryHandler := handlers.NewCategoryHandler(categoryService)

		app := fiber.New()
		app.Delete("/category/:id",categoryHandler.Delete)

		req := httptest.NewRequest("DELETE", "/category/1",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Category deleted successfully")
	})
	t.Run("Error to delete",func(t *testing.T) {

		categoryService := services.NewCategoryServiceMock()
		categoryService.On("DeleteCategory",uint(1)).Return(errors.New("Error to delete category"))

		categoryHandler := handlers.NewCategoryHandler(categoryService)

		app := fiber.New()
		app.Delete("/category/:id",categoryHandler.Delete)

		req := httptest.NewRequest("DELETE", "/category/1",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Error to delete category")
	})
}