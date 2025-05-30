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
	GetAllProducts() ([]models.Product, error)
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

	if product.Price < 0.0 {
		return errors.New("Product price cannot be negative")
	}

	if product.Stock < 0 {
		return errors.New("Product stock cannot be negative")
	}

	if product.Image == "" {
		return errors.New("Product image is required")
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

	if product.Price < 0.0 {
		return errors.New("Product price cannot be negative")
	}

	if product.Stock < 0 {
		return errors.New("Product stock cannot be negative")
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

	existingProduct.Name = product.Name
	existingProduct.Description = product.Description
	existingProduct.Price = product.Price
	existingProduct.Image = product.Image
	existingProduct.Stock = product.Stock
	existingProduct.IsFeatured = product.IsFeatured
	existingProduct.IsOnSale = product.IsOnSale
	existingProduct.CategoryID = product.CategoryID

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

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	products, err := s.productRepo.FindAll()
	if err != nil {
		return nil, errors.New("Error retrieving products")
	}
	return products, nil
}

