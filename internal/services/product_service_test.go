package services_test

import (
	"errors"
	"testing"

	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	repositories "github.com/Beluga-Whale/ecommerce-api/internal/repositories/mocks"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestProduct(t *testing.T){
	t.Run("CreateProduct Success", func(t *testing.T) {
		product := &models.Product{
			Name: "CPU",
			Description:"Test",
			Price: 100,
			Image: "https://www.google.com",
			Stock: 10,
			IsFeatured: false,
			IsOnSale: false,
			CategoryID: 3,
		}

		category := &models.Category{
			Name: "IT",
			Slug: "it",
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByID",product.CategoryID).Return(category,nil)

		productRepo.On("Create",product).Return(nil)

		productService := services.NewProductService(productRepo,categoryRepo)

		err := productService.CreateProduct(product)

		assert.NoError(t,err)

		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)
		
	})

	t.Run("Provide product name , description", func(t *testing.T) {
		product := &models.Product{
			Name: "",
			Description:"",
			Price: 100,
			Image: "https://www.google.com",
			Stock: 10,
			IsFeatured: false,
			IsOnSale: false,
			CategoryID: 3,
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo,categoryRepo)

		err := productService.CreateProduct(product)

		assert.EqualError(t,err,"Please provide product name and description")

	})

	t.Run("Product price cannot be negative", func(t *testing.T) {
		product := &models.Product{
			Name: "CPU",
			Description:"Test",
			Price: -100,
			Image: "https://www.google.com",
			Stock: 10,
			IsFeatured: false,
			IsOnSale: false,
			CategoryID: 3,
		}
		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()
		productService := services.NewProductService(productRepo,categoryRepo)

		err := productService.CreateProduct(product)

		assert.EqualError(t,err,"Product price cannot be negative")

	})

	t.Run("Product price cannot be negative", func(t *testing.T) {
		product := &models.Product{
			Name: "CPU",
			Description:"Test",
			Price: 100,
			Image: "https://www.google.com",
			Stock: -1,
			IsFeatured: false,
			IsOnSale: false,
			CategoryID: 3,
		}
		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()
		productService := services.NewProductService(productRepo,categoryRepo)

		err := productService.CreateProduct(product)

		assert.EqualError(t,err,"Product stock cannot be negative")

	})

	t.Run("Product image is required", func(t *testing.T) {
		product := &models.Product{
			Name: "CPU",
			Description:"Test",
			Price: 100,
			Image: "",
			Stock: 1,
			IsFeatured: false,
			IsOnSale: false,
			CategoryID: 3,
		}
		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()
		productService := services.NewProductService(productRepo,categoryRepo)

		err := productService.CreateProduct(product)

		assert.EqualError(t,err,"Product image is required")

	})

	t.Run("Error finding category", func(t *testing.T) {
		product := &models.Product{
			Name:        "CPU",
			Description: "Test",
			Price:       100,
			Image:       "https://www.google.com",
			Stock:       1,
			IsFeatured:  false,
			IsOnSale:    false,
			CategoryID:  3,
		}
		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()
		productService := services.NewProductService(productRepo, categoryRepo)

		categoryRepo.On("FindByID", product.CategoryID).Return(nil, errors.New("Error finding"))

		err := productService.CreateProduct(product)

		assert.EqualError(t, err, "Error finding category: Error finding")
		categoryRepo.AssertExpectations(t)
	})

	t.Run("Category not found", func(t *testing.T) {
		product := &models.Product{
			Name: "CPU",
			Description:"Test",
			Price: 100,
			Image: "https://www.google.com",
			Stock: 10,
			IsFeatured: false,
			IsOnSale: false,
			CategoryID: 3,
		}


		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByID",product.CategoryID).Return(nil,nil)


		productService := services.NewProductService(productRepo,categoryRepo)

		err := productService.CreateProduct(product)

		assert.EqualError(t,err,"Category not found")
	})

	t.Run("Error creating product", func(t *testing.T) {
		product := &models.Product{
			Name: "CPU",
			Description:"Test",
			Price: 100,
			Image: "https://www.google.com",
			Stock: 10,
			IsFeatured: false,
			IsOnSale: false,
			CategoryID: 3,
		}

		category := &models.Category{
			Name: "IT",
			Slug: "it",
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByID",product.CategoryID).Return(category,nil)

		productRepo.On("Create",product).Return(errors.New("error to create product"))

		productService := services.NewProductService(productRepo,categoryRepo)

		err := productService.CreateProduct(product)

		assert.EqualError(t,err,"Error creating product: error to create product")

		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)
	})
}

func TestUpdateProduct(t *testing.T){
	t.Run("Update Product Success", func(t *testing.T) {
		productID := uint(1)

		updateInput := &models.Product{
			Name:        "Updated CPU",
			Description: "Updated Description",
			Price:       150,
			Image:       "https://example.com/cpu.png",
			Stock:       20,
			IsFeatured:  true,
			IsOnSale:    true,
			CategoryID:  2,
		}

		// Mock existing product in DB
		existingProduct := &models.Product{
			Model:       gorm.Model{ID: productID},
			Name:        "Old CPU",
			Description: "Old Description",
			Price:       100,
			Image:       "https://example.com/old.png",
			Stock:       10,
			IsFeatured:  false,
			IsOnSale:    false,
			CategoryID:  1,
		}

		// Mock existing category
		existingCategory := &models.Category{
			Model: gorm.Model{ID: 2},
			Name:  "Electronics",
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		// Setup mocks
		productRepo.On("FindByID", productID).Return(existingProduct, nil)
		categoryRepo.On("FindByID", updateInput.CategoryID).Return(existingCategory, nil)
		productRepo.On("Update", mock.AnythingOfType("*models.Product")).Return(nil)

		// Call the function
		err := productService.UpdateProduct(productID, updateInput)

		// Assertion
		assert.NoError(t, err)
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)
	})

	t.Run("Please provide product name and description",func(t *testing.T) {
		productID := uint(1)

		updateInput := &models.Product{
			Name:        "",
			Description: "",
			Price:       150,
			Image:       "https://example.com/cpu.png",
			Stock:       20,
			IsFeatured:  true,
			IsOnSale:    true,
			CategoryID:  2,
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.UpdateProduct(productID, updateInput)

		assert.EqualError(t, err,"Please provide product name and description")
	})

	t.Run("Product price cannot be negative",func(t *testing.T) {
		productID := uint(1)

		updateInput := &models.Product{
			Name:        "Updated CPU",
			Description: "Updated Description",
			Price:       -100,
			Image:       "https://example.com/cpu.png",
			Stock:       20,
			IsFeatured:  true,
			IsOnSale:    true,
			CategoryID:  2,
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.UpdateProduct(productID, updateInput)

		assert.EqualError(t, err,"Product price cannot be negative")
	})

	t.Run("Product price cannot be negative",func(t *testing.T) {
		productID := uint(1)

		updateInput := &models.Product{
			Name:        "Updated CPU",
			Description: "Updated Description",
			Price:       100,
			Image:       "https://example.com/cpu.png",
			Stock:       -1,
			IsFeatured:  true,
			IsOnSale:    true,
			CategoryID:  2,
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.UpdateProduct(productID, updateInput)

		assert.EqualError(t, err,"Product stock cannot be negative")
	})

	t.Run("Error finding product", func(t *testing.T) {
		productID := uint(1)

		updateInput := &models.Product{
			Name:        "Updated CPU",
			Description: "Updated Description",
			Price:       150,
			Image:       "https://example.com/cpu.png",
			Stock:       20,
			IsFeatured:  true,
			IsOnSale:    true,
			CategoryID:  2,
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		productRepo.On("FindByID", productID).Return(nil, errors.New("Error finding product"))

		err := productService.UpdateProduct(productID, updateInput)

		assert.EqualError(t, err, "Error finding product")

	})

	t.Run("Product not found", func(t *testing.T) {
		productID := uint(1)

		updateInput := &models.Product{
			Name:        "Updated CPU",
			Description: "Updated Description",
			Price:       150,
			Image:       "https://example.com/cpu.png",
			Stock:       20,
			IsFeatured:  true,
			IsOnSale:    true,
			CategoryID:  2,
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		productRepo.On("FindByID", productID).Return(nil, nil)

		err := productService.UpdateProduct(productID, updateInput)

		assert.EqualError(t, err, "Product not found")

	})
	
	t.Run("Error finding category", func(t *testing.T) {
		productID := uint(1)

		updateInput := &models.Product{
			Name:        "Updated CPU",
			Description: "Updated Description",
			Price:       150,
			Image:       "https://example.com/cpu.png",
			Stock:       20,
			IsFeatured:  true,
			IsOnSale:    true,
			CategoryID:  2,
		}

		// Mock existing product in DB
		existingProduct := &models.Product{
			Model:       gorm.Model{ID: productID},
			Name:        "Old CPU",
			Description: "Old Description",
			Price:       100,
			Image:       "https://example.com/old.png",
			Stock:       10,
			IsFeatured:  false,
			IsOnSale:    false,
			CategoryID:  1,
		}


		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		productRepo.On("FindByID", productID).Return(existingProduct, nil)
		categoryRepo.On("FindByID", updateInput.CategoryID).Return(nil, errors.New("can't to find"))

		err := productService.UpdateProduct(productID, updateInput)

		assert.EqualError(t, err,"Error finding category: can't to find")

	})

	t.Run("Category not found", func(t *testing.T) {
		productID := uint(1)

		updateInput := &models.Product{
			Name:        "Updated CPU",
			Description: "Updated Description",
			Price:       150,
			Image:       "https://example.com/cpu.png",
			Stock:       20,
			IsFeatured:  true,
			IsOnSale:    true,
			CategoryID:  2,
		}

		// Mock existing product in DB
		existingProduct := &models.Product{
			Model:       gorm.Model{ID: productID},
			Name:        "Old CPU",
			Description: "Old Description",
			Price:       100,
			Image:       "https://example.com/old.png",
			Stock:       10,
			IsFeatured:  false,
			IsOnSale:    false,
			CategoryID:  1,
		}


		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		productRepo.On("FindByID", productID).Return(existingProduct, nil)
		categoryRepo.On("FindByID", updateInput.CategoryID).Return(nil,nil)

		err := productService.UpdateProduct(productID, updateInput)

		assert.EqualError(t, err,"Category not found")
	})

	t.Run("Error updating product", func(t *testing.T) {
		productID := uint(1)

		updateInput := &models.Product{
			Name:        "Updated CPU",
			Description: "Updated Description",
			Price:       150,
			Image:       "https://example.com/cpu.png",
			Stock:       20,
			IsFeatured:  true,
			IsOnSale:    true,
			CategoryID:  2,
		}

		// Mock existing product in DB
		existingProduct := &models.Product{
			Model:       gorm.Model{ID: productID},
			Name:        "Old CPU",
			Description: "Old Description",
			Price:       100,
			Image:       "https://example.com/old.png",
			Stock:       10,
			IsFeatured:  false,
			IsOnSale:    false,
			CategoryID:  1,
		}

		// Mock existing category
		existingCategory := &models.Category{
			Model: gorm.Model{ID: 2},
			Name:  "Electronics",
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		// Setup mocks
		productRepo.On("FindByID", productID).Return(existingProduct, nil)
		categoryRepo.On("FindByID", updateInput.CategoryID).Return(existingCategory, nil)
		productRepo.On("Update", mock.AnythingOfType("*models.Product")).Return(errors.New("Error updating product"))

		// Call the function
		err := productService.UpdateProduct(productID, updateInput)

		// Assertion
		assert.EqualError(t, err,"Error updating product")
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)
	})
}

func TestDeleteProduct(t  *testing.T) {
	t.Run("Delete Product Success",func(t *testing.T) {
		productID := uint(1)

		existingProduct := &models.Product{
			Model: gorm.Model{ID: productID},
			Name: "CPU",
			Description:"Test",
			Price: 100,
			Image: "https://www.google.com",
			Stock: 10,
			IsFeatured: false,
			IsOnSale: false,
			CategoryID: 3,
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID",productID).Return(existingProduct,nil)
		productRepo.On("Delete",productID).Return(nil)

		productService := services.NewProductService(productRepo,categoryRepo)

		err := productService.DeleteProduct(productID)

		assert.NoError(t,err)
		productRepo.AssertExpectations(t)
		
	})
	t.Run("Error finding product",func(t *testing.T) {
		productID := uint(1)
		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID",productID).Return(nil,errors.New("Error finding product"))

		productService := services.NewProductService(productRepo,categoryRepo)

		err := productService.DeleteProduct(productID)

		assert.EqualError(t,err,"Error finding product")
	})
	t.Run("Product not found",func(t *testing.T) {
		productID := uint(1)

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID",productID).Return(nil,nil)

		productService := services.NewProductService(productRepo,categoryRepo)

		err := productService.DeleteProduct(productID)

		assert.EqualError(t,err,"Product not found")
		productRepo.AssertExpectations(t)
		
	})
	t.Run("Error deleting product",func(t *testing.T) {
		productID := uint(1)

		existingProduct := &models.Product{
			Model: gorm.Model{ID: productID},
			Name: "CPU",
			Description:"Test",
			Price: 100,
			Image: "https://www.google.com",
			Stock: 10,
			IsFeatured: false,
			IsOnSale: false,
			CategoryID: 3,
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID",productID).Return(existingProduct,nil)
		productRepo.On("Delete",productID).Return(errors.New("Error deleting product"))

		productService := services.NewProductService(productRepo,categoryRepo)

		err := productService.DeleteProduct(productID)

		assert.EqualError(t,err,"Error deleting product")
		
	})
}

func TestGetProductByID(t *testing.T) {
	t.Run("GetProductByID Success",func(t *testing.T) {
		productID := uint(1)

		existingProduct := &models.Product{
			Model: gorm.Model{ID: productID},
			Name: "CPU",
			Description:"Test",
			Price: 100,
			Image: "https://www.google.com",
			Stock: 10,
			IsFeatured: false,
			IsOnSale: false,
			CategoryID: 3,
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID",productID).Return(existingProduct,nil)

		productService := services.NewProductService(productRepo,categoryRepo)

		_ ,err := productService.GetProductByID(productID)

		assert.NoError(t, err)

		productRepo.AssertExpectations(t)

	})

	t.Run("Error finding product",func(t *testing.T) {
		productID := uint(1)

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID",productID).Return(nil,errors.New("Error finding product"))

		productService := services.NewProductService(productRepo,categoryRepo)

		_ ,err := productService.GetProductByID(productID)

		assert.EqualError(t, err,"Error finding product")

		productRepo.AssertExpectations(t)

	})

	t.Run("Error finding product",func(t *testing.T) {
		productID := uint(1)

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID",productID).Return(nil,nil)

		productService := services.NewProductService(productRepo,categoryRepo)

		_ ,err := productService.GetProductByID(productID)

		assert.EqualError(t, err,"Product not found")

		productRepo.AssertExpectations(t)

	})
}

func TestGetAllProducts(t *testing.T) {
	t.Run("GetAllProducts Success",func(t *testing.T) {
		products := []models.Product{
			{
				Name: "CPU",
				Description:"Test",
				Price: 100,
				Image: "https://www.google.com",
				Stock: 10,
				IsFeatured: false,
				IsOnSale: false,
				CategoryID: 3,
			},
			{
				Name: "CPU2",
				Description:"Test2",
				Price: 100,
				Image: "https://www.google.com",
				Stock: 10,
				IsFeatured: false,
				IsOnSale: false,
				CategoryID: 3,
			},
		}

		productRepo := repositories.NewProductRepositoryMock()

		productRepo.On("FindAll",uint(1),uint(10)).Return(products,int64(1),nil)

		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo,categoryRepo)

		_,_,err := productService.GetAllProducts(1,10)

		assert.NoError(t,err)

		productRepo.AssertExpectations(t)

	})

	t.Run("Error retrieving products",func(t *testing.T) {	

		productRepo := repositories.NewProductRepositoryMock()

		productRepo.On("FindAll",uint(1),uint(10)).Return(nil,int64(0),errors.New("Error retrieving products"))

		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo,categoryRepo)

		_,_,err := productService.GetAllProducts(1,10)

		assert.EqualError(t,err,"Error retrieving products")
	})
}