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
	t.Run("Crate Category Success",func(t *testing.T) {
		reqMock := &models.Category{
			Model: gorm.Model{ID: 1},
			Name: "TEST A",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByName",reqMock.Name).Return(nil,nil)
		categoryRepo.On("Create",reqMock).Return(nil)

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.CreateCategory(reqMock)

		assert.NoError(t,err)

		categoryRepo.AssertExpectations(t)

	})

	t.Run("Category Name is nil",func(t *testing.T) {
		reqMock := &models.Category{
			Model: gorm.Model{ID: 1},
			Name: "",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.CreateCategory(reqMock)

		assert.EqualError(t,err,"Category name cannot be empty")

		categoryRepo.AssertExpectations(t)

	})

	t.Run("Error to find Cate By Name",func(t *testing.T) {
		reqMock := &models.Category{
			Model: gorm.Model{ID: 1},
			Name: "TEST A",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByName",reqMock.Name).Return(nil,errors.New("Error checking for existing category"))

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.CreateCategory(reqMock)

		assert.EqualError(t,err,"Error checking for existing category")

		categoryRepo.AssertExpectations(t)

	})

	t.Run("Already existingCategory",func(t *testing.T) {
		reqMock := &models.Category{
			Model: gorm.Model{ID: 1},
			Name: "TEST A",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByName",reqMock.Name).Return(reqMock,nil)

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.CreateCategory(reqMock)

		assert.EqualError(t,err,"Category already exists")

		categoryRepo.AssertExpectations(t)

	})

	t.Run("Error To Create Category",func(t *testing.T) {
		reqMock := &models.Category{
			Model: gorm.Model{ID: 1},
			Name: "TEST A",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByName",reqMock.Name).Return(nil,nil)
		categoryRepo.On("Create",reqMock).Return(errors.New("Error creating category"))

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.CreateCategory(reqMock)

		assert.EqualError(t,err,"Error creating category")

		categoryRepo.AssertExpectations(t)

	})
}

func TestUpdateCategory(t *testing.T) {
	t.Run("Update Category",func(t *testing.T) {
		id := uint(1)
		reqMock := &models.Category{
			Model: gorm.Model{ID: 1},
			Name: "TEST B",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByID",reqMock.Model.ID).Return(reqMock,nil)
		categoryRepo.On("Update",reqMock).Return(nil)

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.UpdateCategory(id,reqMock)

		assert.NoError(t,err)

		categoryRepo.AssertExpectations(t)

	})

	t.Run("Category Name is Empty",func(t *testing.T) {
		id := uint(1)
		reqMock := &models.Category{
			Model: gorm.Model{ID: 1},
			Name: "",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()


		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.UpdateCategory(id,reqMock)

		assert.EqualError(t,err,"Category name cannot be empty")

		categoryRepo.AssertExpectations(t)

	})

	t.Run("Error to find Category By ID",func(t *testing.T) {
		id := uint(1)
		reqMock := &models.Category{
			Model: gorm.Model{ID: 1},
			Name: "TEST B",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByID",reqMock.Model.ID).Return(nil,errors.New("Error finding category"))

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.UpdateCategory(id,reqMock)

		assert.EqualError(t,err,"Error finding category")

		categoryRepo.AssertExpectations(t)

	})

	t.Run("Already Existing Category",func(t *testing.T) {
		id := uint(1)
		reqMock := &models.Category{
			Model: gorm.Model{ID: 1},
			Name: "TEST B",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByID",reqMock.Model.ID).Return(nil,nil)

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.UpdateCategory(id,reqMock)

		assert.EqualError(t,err,"Category not found")

		categoryRepo.AssertExpectations(t)

	})

	t.Run("Error To Update Category",func(t *testing.T) {
		id := uint(1)
		reqMock := &models.Category{
			Model: gorm.Model{ID: 1},
			Name: "TEST B",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByID",reqMock.Model.ID).Return(reqMock,nil)
		categoryRepo.On("Update",reqMock).Return(errors.New("Error updating category"))

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.UpdateCategory(id,reqMock)

		assert.EqualError(t,err,"Error updating category")

		categoryRepo.AssertExpectations(t)

	})
}

func TestDeleteCategory(t *testing.T) {
	t.Run("DeleteCategory Success",func(t *testing.T) {
		id := uint(1)
		reqMock := &models.Category{
			Model: gorm.Model{ID: 1},
			Name: "TEST B",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByID",reqMock.Model.ID).Return(reqMock,nil)
		categoryRepo.On("Delete",id).Return(nil)

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.DeleteCategory(id)

		assert.NoError(t,err)

		categoryRepo.AssertExpectations(t)

	})

	t.Run("Error to FindByID",func(t *testing.T) {
		id := uint(1)
		reqMock := &models.Category{
			Model: gorm.Model{ID: 1},
			Name: "TEST B",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByID",reqMock.Model.ID).Return(nil,errors.New("Error finding category"))

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.DeleteCategory(id)

		assert.EqualError(t,err,"Error finding category")

		categoryRepo.AssertExpectations(t)

	})

	t.Run("Not have category id ",func(t *testing.T) {
		id := uint(1)
		reqMock := &models.Category{
			Model: gorm.Model{ID: 1},
			Name: "TEST B",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByID",reqMock.Model.ID).Return(nil,nil)

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.DeleteCategory(id)

		assert.EqualError(t,err,"Category not found")

		categoryRepo.AssertExpectations(t)

	})

	t.Run("Error to delete category",func(t *testing.T) {
		id := uint(1)
		reqMock := &models.Category{
			Model: gorm.Model{ID: 1},
			Name: "TEST B",
		}

		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByID",reqMock.Model.ID).Return(reqMock,nil)
		categoryRepo.On("Delete",id).Return(errors.New("Error deleting category"))

		categoryService := services.NewCategoryService(categoryRepo)

		err := categoryService.DeleteCategory(id)

		assert.EqualError(t,err,"Error deleting category")

		categoryRepo.AssertExpectations(t)

	})
}

func TestGetAllCategories(t *testing.T) {
	t.Run("GetAllCategories Success",func(t *testing.T) {
		categorys :=[]models.Category{
		{	
			Model: gorm.Model{ID: 1},
			Name: "TEST A",
		},

		{	
			Model: gorm.Model{ID: 2},
			Name: "TEST B",
		},

		}

		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindAll").Return(categorys,nil)

		categoryService := services.NewCategoryService(categoryRepo)

		categoryList,err := categoryService.GetAllCategories()

		assert.NoError(t,err)

		assert.Equal(t,categoryList[0].Name,"TEST A")
		assert.Equal(t,categoryList[1].Name,"TEST B")

		categoryRepo.AssertExpectations(t)
	})

	t.Run("Error To GetAllCategories",func(t *testing.T) {
		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindAll").Return(nil,errors.New("Error retrieving categories"))

		categoryService := services.NewCategoryService(categoryRepo)

		categoryList,err := categoryService.GetAllCategories()

		assert.EqualError(t,err,"Error retrieving categories")
		assert.Nil(t,categoryList)

		categoryRepo.AssertExpectations(t)
	})
}