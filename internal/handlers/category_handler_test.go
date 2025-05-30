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

func TestGetAll(t *testing.T) {
	t.Run("Get all categories success", func(t *testing.T) {
		category := []models.Category{
			{
				Name: "Electronics",
			},
			{
				Name: "Books",
			},
		}

		app := fiber.New()
		categoryService := services.NewCategoryServiceMock()
		categoryHandler := handlers.NewCategoryHandler(categoryService)
		app.Get("/category", categoryHandler.GetAll)

		categoryService.On("GetAllCategories").Return(category,nil)

		req := httptest.NewRequest("GET","/category", nil)

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Get all categories successfully")

	})

	t.Run("Failed to fetch categories", func(t *testing.T) {
		app := fiber.New()
		categoryService := services.NewCategoryServiceMock()
		categoryHandler := handlers.NewCategoryHandler(categoryService)
		app.Get("/category", categoryHandler.GetAll)

		categoryService.On("GetAllCategories").Return(nil,errors.New("Failed to fetch categories"))

		req := httptest.NewRequest("GET","/category", nil)

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Failed to fetch categories")

	})
}
func TestCreate(t *testing.T){
	t.Run("Create category success", func(t *testing.T) {
		app := fiber.New()
		categoryService := services.NewCategoryServiceMock()
		categoryHandler := handlers.NewCategoryHandler(categoryService)
		app.Post("/category", categoryHandler.Create)

		category := &models.Category{
			Name: "Electronics",
		}

		categoryService.On("CreateCategory", category).Return(nil)

		reqBody := []byte(`{"name":"Electronics"}`)

		req := httptest.NewRequest("POST","/category", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Category created successfully")
	})
	t.Run("Invalid request body", func(t *testing.T) {
		app := fiber.New()
		categoryService := services.NewCategoryServiceMock()
		categoryHandler := handlers.NewCategoryHandler(categoryService)
		app.Post("/category", categoryHandler.Create)

		reqBody := []byte(`{"name":"Electronics",}`)

		req := httptest.NewRequest("POST","/category", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Invalid request body")
	})
	t.Run("Validation error", func(t *testing.T) {
		app := fiber.New()
		categoryService := services.NewCategoryServiceMock()
		categoryHandler := handlers.NewCategoryHandler(categoryService)
		app.Post("/category", categoryHandler.Create)

		tests := []struct{
			name string
			reqBody []byte
			expect []string
		}{
			{
				name: "Name is required",
				reqBody: []byte(`{"name":""}`),
				expect: []string{"Name is required"},
			},
			{
				name: "Name is min",
				reqBody: []byte(`{"name":"t"}`),
				expect: []string{"Name is min"},
			},
		}

		for _,tc := range tests {
			t.Run(tc.name,func(t *testing.T) {
				req := httptest.NewRequest("POST","/category", bytes.NewReader(tc.reqBody))

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

	t.Run("Error create category", func(t *testing.T) {
		app := fiber.New()
		categoryService := services.NewCategoryServiceMock()
		categoryHandler := handlers.NewCategoryHandler(categoryService)
		app.Post("/category", categoryHandler.Create)

		category := &models.Category{
			Name: "Electronics",
		}

		categoryService.On("CreateCategory", category).Return(errors.New("Error creating category"))

		reqBody := []byte(`{"name":"Electronics"}`)

		req := httptest.NewRequest("POST","/category", bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t,string(body),"Error creating category")
	})
}
func TestUpdate(t *testing.T){
	t.Run("Update category success", func(t *testing.T) {
		app := fiber.New()
		categoryService := services.NewCategoryServiceMock()
		categoryHandler := handlers.NewCategoryHandler(categoryService)
		app.Put("/category/:id",categoryHandler.Update)

		category := &models.Category{
			Name: "Electronics",
		}

		categoryService.On("UpdateCategory", uint(1), category).Return(nil)

		reqBody := []byte(`{"name":"Electronics"}`)

		req := httptest.NewRequest("PUT", "/category/1",bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Category update successfully")
	})

	t.Run("Invalid request body", func(t *testing.T) {
		app := fiber.New()
		categoryService := services.NewCategoryServiceMock()
		categoryHandler := handlers.NewCategoryHandler(categoryService)
		app.Put("/category/:id",categoryHandler.Update)

		reqBody := []byte(`{"name":"Electronics",}`)

		req := httptest.NewRequest("PUT", "/category/1",bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request bod")
	})

	t.Run("Validation error", func(t *testing.T) {
		app := fiber.New()
		categoryService := services.NewCategoryServiceMock()
		categoryHandler := handlers.NewCategoryHandler(categoryService)
		app.Put("/category/:id",categoryHandler.Update)

		test := []struct {
			name string
			reqBody []byte
			expect []string
		}{
			{
				name: "Name is required",
				reqBody: []byte(`{"name":""}`),
				expect: []string{"Name is required"},
			},
			{
				name: "Name is min",
				reqBody: []byte(`{"name":"t"}`),
				expect: []string{"Name is min"},
			},
		}

		for _, tc := range test{
			t.Run(tc.name,func(t *testing.T) {
				reqBody := tc.reqBody
				req := httptest.NewRequest("PUT", "/category/1",bytes.NewReader(reqBody))

				req.Header.Set("Content-Type", "application/json")
				res, err := app.Test(req)

				assert.NoError(t, err)
				assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
				body, _ := io.ReadAll(res.Body)

				for _,expect := range tc.expect{
					assert.Contains(t, string(body),expect)
				}	
			})
		}
	})

	t.Run("Error update category", func(t *testing.T) {
		app := fiber.New()
		categoryService := services.NewCategoryServiceMock()
		categoryHandler := handlers.NewCategoryHandler(categoryService)
		app.Put("/category/:id",categoryHandler.Update)

		category := &models.Category{
			Name: "Electronics",
		}

		categoryService.On("UpdateCategory", uint(1), category).Return(errors.New("Error updating category"))

		reqBody := []byte(`{"name":"Electronics"}`)

		req := httptest.NewRequest("PUT", "/category/1",bytes.NewReader(reqBody))

		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Error updating category")
	})
}
func TestDelete(t *testing.T){
	t.Run("Delete category success", func(t *testing.T) {
		app := fiber.New()
		categoryService := services.NewCategoryServiceMock()
		categoryHandler := handlers.NewCategoryHandler(categoryService)
		app.Delete("/category/:id", categoryHandler.Delete)

		categoryService.On("DeleteCategory", uint(1)).Return(nil)

		req := httptest.NewRequest("DELETE", "/category/1", nil)

		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Category deleted successfully")
	})

	t.Run("Error Delete Category", func(t *testing.T) {
		app := fiber.New()
		categoryService := services.NewCategoryServiceMock()
		categoryHandler := handlers.NewCategoryHandler(categoryService)
		app.Delete("/category/:id", categoryHandler.Delete)

		categoryService.On("DeleteCategory", uint(1)).Return(errors.New("Error deleting category"))

		req := httptest.NewRequest("DELETE", "/category/1", nil)

		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Error deleting category")
	})

}