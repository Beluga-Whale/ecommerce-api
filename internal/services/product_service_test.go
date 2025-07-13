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

func TestCreateProduct(t *testing.T) {
	t.Run("Create Success", func(t *testing.T) {
		productCategoryID := uint(1)
		salePrice := 50.0
		product := &models.Product{
			Name: "T-shirt",
			Title: "Title",
			Description: "Test Description",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByID", productCategoryID ).Return(&models.Category{
			Model: gorm.Model{ID: 1},
			Name: "Test Category",
		}, nil)

		productRepo.On("Create", product).Return(nil)

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.CreateProduct(product)

		assert.NoError(t, err)

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Please provide product name and description", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "",
			Title: "Title",
			Description: "",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.CreateProduct(product)

		assert.EqualError(t, err,"Please provide product name and description")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Please provide title product", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "TEST",
			Title: "",
			Description: "ES",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.CreateProduct(product)

		assert.EqualError(t, err,"Please provide title product")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("SalePrice is nil", func(t *testing.T) {
		product := &models.Product{
			Name: "TEST",
			Title: "Title",
			Description: "ES",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   nil,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.CreateProduct(product)

		assert.EqualError(t, err,"Sale price must be greater than 0")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Is on sale FALSE but sale price more than 0", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "TEST",
			Title: "TEST",
			Description: "ES",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    false,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.CreateProduct(product)

		assert.EqualError(t, err,"You can should is in sale true")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})
		

	t.Run("Loop Product variant", func(t *testing.T) {
		// NOTE - สร้าง array แต่ละ test case
		sale := 30.0	
		type testCase struct {
			name  string
			variant models.ProductVariant
			salePrice *float64
			expected string
		}

		cases := []testCase {
			{
				name: "Stock negative",
				variant: models.ProductVariant{
					Size: "S",
					Stock: -10,
					Price: 1000,
				},
				salePrice: &sale,
				expected: "Product stock cannot be negative",
			},
			{
				name: "Price negative",
				variant: models.ProductVariant{
					Size: "S",
					Stock: 10,
					Price: -1000,
				},
				salePrice: &sale,
				expected: "Product price cannot be negative",
			},
			{
				name: "Final Price less than 0",
				variant: models.ProductVariant{
					Size: "S",
					Stock: 10,
					Price: 20,
				},
				salePrice: &sale,
				expected: "Final price of variant 'S' cannot be negative",
			},
			{
				name: "Final Price less than 0",
				variant: models.ProductVariant{
					Size: "S",
					Stock: 10,
					Price: 30,
				},
				salePrice: &sale,
				expected: "Sale price must be less than variant price for size 'S'",
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		for _, item := range cases {
			t.Run(item.name,func(t *testing.T) {
				product := &models.Product{
						Name: "TEST",
						Title: "TEST",
						Description: "ES",
						Images: []models.ProductImage{
							{
								URL: "test",
							},{
								URL: "test",
							},{
								URL: "test",
							},
						},
						IsOnSale:    true,
						SalePrice:   item.salePrice,
						CategoryID:  1,
						Variants: []models.ProductVariant{
							item.variant,
						},
							}
				err := productService.CreateProduct(product)

			    assert.Error(t,err)
				assert.EqualError(t, err,item.expected)

				// NOTE - เช็คว่ามีการ Call function ไหม
				productRepo.AssertExpectations(t)
				categoryRepo.AssertExpectations(t)
			})
		}
	})

	t.Run("Image less than 3 ", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "T-shirt",
			Title: "Title",
			Description: "Test Description",
			Images: []models.ProductImage{
				{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.CreateProduct(product)

		assert.Error(t, err)
		assert.EqualError(t,err,"You must upload exactly 3 product images")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Can't to find category",func(t *testing.T) {
		productCategoryID := uint(1)
		salePrice := 50.0
		product := &models.Product{
			Name: "T-shirt",
			Title: "Title",
			Description: "Test Description",
			Images: []models.ProductImage{
				{
					URL: "test",
				},
				{
					URL: "test",
				},
				{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		categoryRepo.On("FindByID", productCategoryID ).Return(&models.Category{
			Name: "Test Category",
		}, errors.New("Error finding category"))


		err := productService.CreateProduct(product)

		assert.Error(t, err)
		assert.EqualError(t,err,"Error finding category: Error finding category")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)
	})

	t.Run("Category is nil", func(t *testing.T) {
		productCategoryID := uint(1)
		salePrice := 50.0
		product := &models.Product{
			Name: "T-shirt",
			Title: "Title",
			Description: "Test Description",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByID", productCategoryID ).Return(nil, nil)


		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.CreateProduct(product)

		assert.Error(t, err)
		assert.EqualError(t,err,"Category not found")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)
	})

	t.Run("Register Success", func(t *testing.T) {
		productCategoryID := uint(1)
		salePrice := 50.0
		product := &models.Product{
			Name: "T-shirt",
			Title: "Title",
			Description: "Test Description",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		categoryRepo.On("FindByID", productCategoryID ).Return(&models.Category{
			Name: "Test Category",
		}, nil)

		productRepo.On("Create", product).Return(errors.New("Error creating product"))

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.CreateProduct(product)

		assert.Error(t, err)
		assert.EqualError(t, err ,"Error creating product: Error creating product")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})
}

func TestUpdateProduct(t *testing.T) {
	t.Run("Update Success", func(t *testing.T) {
		id:= uint(1)
		categoryId := uint(1)
		salePrice := 50.0
		product := &models.Product{
			Name: "T-shirt",
			Title: "Title",
			Description: "Test Description",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID", id).Return(product,nil)

		categoryRepo.On("FindByID", categoryId ).Return(&models.Category{
			Name: "Test Category",
		}, nil)

		productRepo.On("Update", product).Return(nil)

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.UpdateProduct(1,product)

		assert.NoError(t, err)

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Please provide product name and description", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "",
			Title: "Title",
			Description: "",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.UpdateProduct(1,product)

		assert.EqualError(t, err,"Please provide product name and description")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Please provide title product", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "TEST",
			Title: "",
			Description: "ES",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.UpdateProduct(1,product)

		assert.EqualError(t, err,"Please provide title product")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Image is equal 3", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "T-shirt",
			Title: "Title",
			Description: "Test Description",
			Images: []models.ProductImage{
				{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.UpdateProduct(1,product)

		assert.EqualError(t, err,"You must upload exactly 3 product images")

	})

	t.Run("SalePrice is nil", func(t *testing.T) {
		product := &models.Product{
			Name: "TEST",
			Title: "Title",
			Description: "ES",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   nil,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.UpdateProduct(1,product)

		assert.EqualError(t, err,"Sale price must be greater than 0")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Is on sale FALSE but sale price more than 0", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "TEST",
			Title: "TEST",
			Description: "ES",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    false,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.UpdateProduct(1,product)

		assert.EqualError(t, err,"You can should is in sale true")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})
		

	t.Run("Loop Product variant", func(t *testing.T) {
		// NOTE - สร้าง array แต่ละ test case
		sale := 30.0	
		type testCase struct {
			name  string
			variant models.ProductVariant
			salePrice *float64
			expected string
		}

		cases := []testCase {
			{
				name: "Stock negative",
				variant: models.ProductVariant{
					Size: "S",
					Stock: -10,
					Price: 1000,
				},
				salePrice: &sale,
				expected: "Product stock cannot be negative",
			},
			{
				name: "Price negative",
				variant: models.ProductVariant{
					Size: "S",
					Stock: 10,
					Price: -1000,
				},
				salePrice: &sale,
				expected: "Product price cannot be negative",
			},
			{
				name: "Final Price less than 0",
				variant: models.ProductVariant{
					Size: "S",
					Stock: 10,
					Price: 20,
				},
				salePrice: &sale,
				expected: "Final price of variant 'S' cannot be negative",
			},
			{
				name: "Final Price less than 0",
				variant: models.ProductVariant{
					Size: "S",
					Stock: 10,
					Price: 30,
				},
				salePrice: &sale,
				expected: "Sale price must be less than variant price 'S'",
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productService := services.NewProductService(productRepo, categoryRepo)

		for _, item := range cases {
			t.Run(item.name,func(t *testing.T) {
				product := &models.Product{
						Name: "TEST",
						Title: "TEST",
						Description: "ES",
						Images: []models.ProductImage{
							{
								URL: "test",
							},{
								URL: "test",
							},{
								URL: "test",
							},
						},
						IsOnSale:    true,
						SalePrice:   item.salePrice,
						CategoryID:  1,
						Variants: []models.ProductVariant{
							item.variant,
						},
							}
				err := productService.UpdateProduct(1,product)

			    assert.Error(t,err)
				assert.EqualError(t, err,item.expected)

				// NOTE - เช็คว่ามีการ Call function ไหม
				productRepo.AssertExpectations(t)
				categoryRepo.AssertExpectations(t)
			})
		}
	})

	t.Run("Can't to find product by ID", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "T-shirt",
			Title: "Title",
			Description: "Test Description",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID", mock.Anything).Return(nil,errors.New("Error finding product"))

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.UpdateProduct(1,product)

		assert.Error(t, err)
		assert.EqualError(t,err,"Error finding product")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Get Product ByID is nil", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "T-shirt",
			Title: "Title",
			Description: "Test Description",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID", mock.Anything).Return(nil,nil)

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.UpdateProduct(1,product)

		assert.Error(t, err)
		assert.EqualError(t,err,"Product not found")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Error to find category", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "T-shirt",
			Title: "Title",
			Description: "Test Description",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID", mock.Anything).Return(product,nil)

		categoryRepo.On("FindByID", mock.Anything ).Return(nil, errors.New("Error finding category"))

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.UpdateProduct(1,product)

		assert.Error(t, err)
		assert.EqualError(t,err,"Error finding category: Error finding category")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Category is Nil", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "T-shirt",
			Title: "Title",
			Description: "Test Description",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID", mock.Anything).Return(product,nil)

		categoryRepo.On("FindByID", mock.Anything ).Return(nil, nil)

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.UpdateProduct(1,product)

		assert.Error(t, err)
		assert.EqualError(t,err,"Category not found")

	})

	t.Run("Update Success", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "T-shirt",
			Title: "Title",
			Description: "Test Description",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID", mock.Anything).Return(product,nil)

		categoryRepo.On("FindByID", mock.Anything ).Return(&models.Category{
			Name: "Test Category",
		}, nil)

		productRepo.On("Update", mock.Anything).Return(errors.New(" deleting product"))

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.UpdateProduct(1,product)

		assert.Error(t, err)
		assert.EqualError(t,err,"Error updating product")
		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})
}

func TestDeleteProduct(t *testing.T) {
	t.Run("Delete Success", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "T-shirt",
			Title: "Title",
			Description: "Test Description",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID", mock.Anything).Return(product,nil)
		productRepo.On("Delete", uint(1)).Return(nil)

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.DeleteProduct(1)

		assert.NoError(t, err)

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Error to find product by id", func(t *testing.T) {
		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID", mock.Anything).Return(nil,errors.New("Error finding product"))

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.DeleteProduct(1)

		assert.Error(t, err)
		assert.EqualError(t,err,"Error finding product")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Product is Nil", func(t *testing.T) {
		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID", mock.Anything).Return(nil,nil)

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.DeleteProduct(1)

		assert.Error(t, err)
		assert.EqualError(t,err,"Product not found")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Error Delete", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "T-shirt",
			Title: "Title",
			Description: "Test Description",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID", mock.Anything).Return(product,nil)
		productRepo.On("Delete", mock.Anything).Return(errors.New("Error deleting product"))

		productService := services.NewProductService(productRepo, categoryRepo)

		err := productService.DeleteProduct(1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Error deleting product")


		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)
	})
}

func TestGetProductByID(t *testing.T) {
	t.Run("GetProductByID Success", func(t *testing.T) {
		salePrice := 50.0
		product := &models.Product{
			Name: "T-shirt",
			Title: "Title",
			Description: "Test Description",
			Images: []models.ProductImage{
				{
					URL: "test",
				},{
					URL: "test",
				},{
					URL: "test",
				},
			},
			IsOnSale:    true,
			SalePrice:   &salePrice,
			CategoryID:  1,
			Variants: []models.ProductVariant{
				{
					Size:  "M",
					Price: 100.0,
					Stock: 10,
				},
			},
		}

		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID", mock.Anything).Return(product,nil)
		productService := services.NewProductService(productRepo, categoryRepo)

		_,err := productService.GetProductByID(1)

		assert.NoError(t, err)

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})
	t.Run("Error to find Get product", func(t *testing.T) {
		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID", mock.Anything).Return(nil,errors.New("Error finding product"))
		productService := services.NewProductService(productRepo, categoryRepo)

		_,err := productService.GetProductByID(1)

		assert.Error(t, err)
		assert.EqualError(t,err,"Error finding product")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

	t.Run("Error to find Get product", func(t *testing.T) {
		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		productRepo.On("FindByID", mock.Anything).Return(nil,nil)
		productService := services.NewProductService(productRepo, categoryRepo)

		_,err := productService.GetProductByID(1)

		assert.Error(t, err)
		assert.EqualError(t,err,"Product not found")

		// NOTE - เช็คว่ามีการ Call function ไหม
		productRepo.AssertExpectations(t)
		categoryRepo.AssertExpectations(t)

	})

}

func TestGetAllProduct(t *testing.T) {
	t.Run("GetAllProducts Success", func(t *testing.T) {
		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()

		mockProducts := []models.Product{
			{Name: "T-shirt"}, {Name: "Hoodie"},
		}
		mockPageTotal := int64(5)

		productRepo.
			On("FindAll", uint(1), uint(10), int64(0), int64(1000), "shirt", []int{1, 2}, []string{"M", "L"}).
			Return(mockProducts, mockPageTotal, nil)

		service := services.NewProductService(productRepo, categoryRepo)

		products, pageTotal, err := service.GetAllProducts(1, 10, 0, 1000, "shirt", []int{1, 2}, []string{"M", "L"})

		assert.NoError(t, err)
		assert.Equal(t, mockProducts, products)
		assert.Equal(t, mockPageTotal, pageTotal)

		productRepo.AssertExpectations(t)
	})

	t.Run("GetAllProducts Success", func(t *testing.T) {
		productRepo := repositories.NewProductRepositoryMock()
		categoryRepo := repositories.NewCategoryRepositoryMock()
		productRepo.
			On("FindAll", uint(1), uint(10), int64(0), int64(1000), "shirt", []int{1, 2}, []string{"M", "L"}).
			Return(nil, nil, errors.New("Error retrieving products"))

		service := services.NewProductService(productRepo, categoryRepo)

		_, _, err := service.GetAllProducts(1, 10, 0, 1000, "shirt", []int{1, 2}, []string{"M", "L"})

		assert.EqualError(t, err, "Error retrieving products")

		productRepo.AssertExpectations(t)
	})
}