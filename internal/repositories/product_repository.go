package repositories

import (
	"errors"

	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type ProductRepositoryInterface interface {
	Create(product *models.Product) error
	FindByID(id uint) (*models.Product, error)
	FindAll() ([]models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) FindByID(id uint) (*models.Product, error){
	var product models.Product

	err := r.db.First(&product,id).Error

	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil,nil
	}

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Find(&products).Error
	return products, err
}

func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}