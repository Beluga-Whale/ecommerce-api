package repositories

import (
	"errors"

	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type CategoryInterface interface {
	Create(category *models.Category) error 
	Update(category *models.Category) error
	Delete(id uint) error
	FindAll() ([]models.Category,error)
	FindByName(name string) (*models.Category,error)
	FindByID(id uint) (*models.Category,error)
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository{
	return &CategoryRepository{db:db}
}

func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

func (r *CategoryRepository) FindAll() ([]models.Category,error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return 	categories, err
}

func (r *CategoryRepository) FindByName(name string) (*models.Category,error) {
	var category models.Category
	err := r.db.Where("name = ?",name).First(&category).Error

	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil,nil
	}

	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) FindByID(id uint) (*models.Category,error) {
	var category models.Category
	err := r.db.First(&category,id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err	
	}

	return 	&category, err
}