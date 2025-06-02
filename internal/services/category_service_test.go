package services_test

import (
	"errors"
	"testing"

	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	repositories "github.com/Beluga-Whale/ecommerce-api/internal/repositories/mocks"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateCategory(t *testing.T) {
	t.Run("Create Category Success", func(t *testing.T) {
		category := &models.Category{
			Name: "Electronics",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()
		categoryRepo.On("FindByName", category.Name).Return(nil, nil)

		categoryRepo.On("Create", category).Return(nil)

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.CreateCategory(category)

		assert.NoError(t, err)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Category name cannot be empty", func(t *testing.T) {
		category := &models.Category{
			Name: "",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()
		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.CreateCategory(category)

		assert.EqualError(t, err,"Category name cannot be empty")

	})

	t.Run("Error checking for existing category", func(t *testing.T) {
		category := &models.Category{
			Name: "Electronics",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()
		categoryRepo.On("FindByName", category.Name).Return(nil, errors.New("Error checking for existing category"))

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.CreateCategory(category)

		assert.EqualError(t,err,"Error checking for existing category")
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Category already exists", func(t *testing.T) {
		category := &models.Category{
			Name: "Electronics",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()
		categoryRepo.On("FindByName", category.Name).Return(category, nil)

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.CreateCategory(category)

		assert.EqualError(t, err,"Category already exists")

	})

	t.Run("Create Category Success", func(t *testing.T) {
		category := &models.Category{
			Name: "Electronics",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()
		categoryRepo.On("FindByName", category.Name).Return(nil, nil)

		categoryRepo.On("Create", category).Return(errors.New("Error creating category"))

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.CreateCategory(category)

		assert.EqualError(t, err, "Error creating category")
		categoryRepo.AssertExpectations(t)

	})
}

func TestUpdateCategory(t *testing.T) {
	t.Run("Update Category Success", func(t *testing.T) {
		category := &models.Category{
			Name: "Electronics",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()
		categoryRepo.On("FindByID", uint(1)).Return(category, nil)
		categoryRepo.On("Update", category).Return(nil)

		categoryService := services.NewCategoryService(categoryRepo)
		err := categoryService.UpdateCategory(1, category)
		assert.NoError(t, err)
		categoryRepo.AssertExpectations(t)
	})
	t.Run("Update Category name cannot be empty", func(t *testing.T) {
		category := &models.Category{
			Name: "",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryService := services.NewCategoryService(categoryRepo)
		err := categoryService.UpdateCategory(1, category)
		assert.EqualError(t, err, "Category name cannot be empty")
		categoryRepo.AssertExpectations(t)
	})
	t.Run("Error finding category", func(t *testing.T) {
		category := &models.Category{
			Name: "Electronics",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()
		categoryRepo.On("FindByID", uint(1)).Return(nil, errors.New("Error finding category"))

		categoryService := services.NewCategoryService(categoryRepo)
		err := categoryService.UpdateCategory(1, category)
		assert.EqualError(t, err, "Error finding category")
		categoryRepo.AssertExpectations(t)
	})
	t.Run("Category not found", func(t *testing.T) {
		category := &models.Category{
			Name: "Electronics",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()
		categoryRepo.On("FindByID", uint(1)).Return(nil, nil)

		categoryService := services.NewCategoryService(categoryRepo)
		err := categoryService.UpdateCategory(1, category)
		assert.EqualError(t, err, "Category not found")
		categoryRepo.AssertExpectations(t)
	})
	t.Run("Error updating category", func(t *testing.T) {
		category := &models.Category{
			Name: "Electronics",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()
		categoryRepo.On("FindByID", uint(1)).Return(category, nil)
		categoryRepo.On("Update", category).Return(errors.New("Error updating category"))

		categoryService := services.NewCategoryService(categoryRepo)
		err := categoryService.UpdateCategory(1, category)
		assert.EqualError(t, err, "Error updating category")
		categoryRepo.AssertExpectations(t)
	})
}

func TestDeleteCategory(t *testing.T) {
	t.Run("Delete Category Success", func(t *testing.T) {
		category := &models.Category{
			Name: "Electronics",
			Model: gorm.Model{
				ID: 1,
			},
		}
		categoryRepo := repositories.NewCategoryRepositoryMock()
		categoryRepo.On("FindByID", uint(1)).Return(category, nil)
		categoryRepo.On("Delete", uint(1)).Return(nil)

		categoryService := services.NewCategoryService(categoryRepo)
		err := categoryService.DeleteCategory(1)

		assert.NoError(t, err)
		categoryRepo.AssertExpectations(t)
	})
	t.Run("Error finding category", func(t *testing.T) {
		categoryRepo := repositories.NewCategoryRepositoryMock()
		categoryRepo.On("FindByID", uint(1)).Return(nil, errors.New("Error finding category"))

		categoryService := services.NewCategoryService(categoryRepo)
		err := categoryService.DeleteCategory(1)

		assert.EqualError(t, err, "Error finding category")
		categoryRepo.AssertExpectations(t)
	})

	t.Run("Category not found", func(t *testing.T) {
		categoryRepo := repositories.NewCategoryRepositoryMock()
		categoryRepo.On("FindByID", uint(1)).Return(nil, nil)

		categoryService := services.NewCategoryService(categoryRepo)
		err := categoryService.DeleteCategory(1)

		assert.EqualError(t, err, "Category not found")
		categoryRepo.AssertExpectations(t)
	})

	t.Run("Error deleting category", func(t *testing.T) {
		category := &models.Category{
			Name: "Electronics",
			Model: gorm.Model{
				ID: 1,
			},
		}
		categoryRepo := repositories.NewCategoryRepositoryMock()
		categoryRepo.On("FindByID", uint(1)).Return(category, nil)
		categoryRepo.On("Delete", uint(1)).Return(errors.New("Error deleting category"))

		categoryService := services.NewCategoryService(categoryRepo)
		err := categoryService.DeleteCategory(1)

		assert.EqualError(t, err, "Error deleting category")
		categoryRepo.AssertExpectations(t)
	})
}

func TestFindAllCategories(t *testing.T) {
	t.Run("Find All Categories Success", func(t *testing.T) {
		categories := []models.Category{
			{Name: "Electronics"},
			{Name: "Books"},
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()
		categoryRepo.On("FindAll").Return(categories, nil)

		categoryService := services.NewCategoryService(categoryRepo)
		result, err := categoryService.GetAllCategories()

		assert.NoError(t, err)
		assert.Equal(t, categories, result)
		categoryRepo.AssertExpectations(t)
	})

	t.Run("Error retrieving categories", func(t *testing.T) {
		categoryRepo := repositories.NewCategoryRepositoryMock()
		categoryRepo.On("FindAll").Return(nil, errors.New("Error retrieving categories"))

		categoryService := services.NewCategoryService(categoryRepo)
		_, err := categoryService.GetAllCategories()

		assert.EqualError(t, err, "Error retrieving categories")
		categoryRepo.AssertExpectations(t)
	})
}