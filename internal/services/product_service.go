package services

import (
	"errors"
	"fmt"

	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/repositories"
)

type ProductServiceInterface interface{
	CreateProduct(product *models.Product) error
	UpdateProduct(id uint, product *models.Product) error
	DeleteProduct(id uint) error
	GetProductByID(id uint) (*models.Product, error) 
	GetAllProducts( page uint, limit uint, minPrice int64, maxPrice int64, searchName string, categoryIDs []int) ([]models.Product, int64,error) 
}

type ProductService struct {
	productRepo repositories.ProductRepositoryInterface
	categoryRepo repositories.CategoryInterface
}

func NewProductService(productRepo repositories.ProductRepositoryInterface, categoryRepo repositories.CategoryInterface) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	if product.Name == "" || product.Description ==""{
		return errors.New("Please provide product name and description")	
	}

	if product.Title == "" {
		return errors.New("Please provide title product")
	}

	if product.IsOnSale {
		if product.SalePrice == nil || *product.SalePrice <= 0.0 {
			return errors.New("Sale price must be greater than 0")
		}
	}

	if product.IsOnSale == false && *product.SalePrice > 0.0 {
		return errors.New("You can should is in sale true")
	}


	for _, v := range product.Variants {
		if v.Stock < 0{
			return errors.New("Product stock cannot be negative")
		}

		if v.Price < 0.0 {
			return errors.New("Product price cannot be negative")
		}

		if product.IsOnSale && product.SalePrice != nil{
			finalPrice := v.Price - *product.SalePrice
			if finalPrice < 0 {
				return fmt.Errorf("Final price of variant '%s' cannot be negative", v.Size)
			}
			if *product.SalePrice >= v.Price {
				return fmt.Errorf("Sale price must be less than variant price for size '%s'", v.Size)
			}
		}
	}


	if len(product.Images) != 3 {
		return errors.New("You must upload exactly 3 product images")
	}

	category,err := s.categoryRepo.FindByID(product.CategoryID)

	if err != nil {
		return fmt.Errorf("Error finding category: %w", err)
	}

	if category == nil {
		return errors.New("Category not found")
	}

	err = s.productRepo.Create(product)
	if err != nil {
		return fmt.Errorf("Error creating product: %w", err)
	}
	return nil
}

func (s *ProductService) UpdateProduct(id uint, product *models.Product) error {
	if product.Name == "" || product.Description == "" {
		return errors.New("Please provide product name and description")
	}

	if product.Title == "" {
		return errors.New("Please provide title product")
	}

	if product.IsOnSale {
		if product.SalePrice == nil || *product.SalePrice <= 0.0 {
			return errors.New("Sale price must be greater than 0")
		}
	}

	if product.IsOnSale == false && *product.SalePrice > 0.0 {
		return errors.New("You can should is in sale true")
	}

	for _, v := range product.Variants {
		if v.Stock < 0{
			return errors.New("Product stock cannot be negative")
		}

		if v.Price < 0.0 {
			return errors.New("Product price cannot be negative")
		}

		if product.IsOnSale && product.SalePrice != nil{
			finalPrice := v.Price - *product.SalePrice
			if finalPrice < 0 {
				return fmt.Errorf("Final price of variant '%s' cannot be negative", v.Size)
			}
			if *product.SalePrice >= v.Price {
				return fmt.Errorf("Sale price must be less than variant price '%s'", v.Size)
			}
		}
	}

	existingProduct, err := s.productRepo.FindByID(id)
	if err != nil {
		return errors.New("Error finding product")
	}

	if existingProduct == nil {
		return errors.New("Product not found")
	}

	category,err := s.categoryRepo.FindByID(product.CategoryID)

	if err != nil {
		return fmt.Errorf("Error finding category: %w", err)
	}

	if category == nil {
		return errors.New("Category not found")
	}

	var variantsUpdate = []models.ProductVariant{}

	for _, v := range product.Variants {
		variantsUpdate = append(variantsUpdate, models.ProductVariant{
			Size: v.Size,
			Stock: v.Stock,
			SKU: v.SKU,
			Price: v.Price,
		})
	}

	var urlUpdate =[]models.ProductImage{}
	for _, i := range product.Images{
		urlUpdate = append(urlUpdate, models.ProductImage{
			URL: i.URL,
		})
	}

	existingProduct.Name = product.Name
	existingProduct.Title = product.Title
	existingProduct.Description = product.Description
	existingProduct.Images = urlUpdate
	existingProduct.IsFeatured = product.IsFeatured
	existingProduct.IsOnSale = product.IsOnSale
	existingProduct.SalePrice = product.SalePrice
	existingProduct.CategoryID = product.CategoryID
	existingProduct.Variants = variantsUpdate

	err = s.productRepo.Update(existingProduct)
	if err != nil {
		return errors.New("Error updating product")
	}
	return nil
}

func (s *ProductService) DeleteProduct(id uint) error {
	existingProduct, err := s.productRepo.FindByID(id)
	if err != nil {
		return errors.New("Error finding product")
	}
	if existingProduct == nil {
		return errors.New("Product not found")
	}
	err = s.productRepo.Delete(id)

	if err != nil {
		return errors.New("Error deleting product")
	}
	return nil
}

func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	existingProduct, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("Error finding product")
	}	
	if existingProduct == nil {
		return nil, errors.New("Product not found")
	}
	return existingProduct, nil
}

func (s *ProductService) GetAllProducts( page uint, limit uint, minPrice int64, maxPrice int64, searchName string, categoryIDs []int) ([]models.Product, int64,error) {
	products,pageTotal, err := s.productRepo.FindAll(page,limit,minPrice,maxPrice,searchName ,categoryIDs)
	if err != nil {
		return nil, 0,errors.New("Error retrieving products")
	}
	return products,pageTotal, nil
}

