package services

import (
	"errors"

	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/repositories"
	"github.com/gosimple/slug"
)

type CategoryServiceInterface interface {
	CreateCategory(category *models.Category) error
	UpdateCategory(id uint, category *models.Category) error
	DeleteCategory(id uint) error
	GetAllCategories() ([]models.Category, error)
}

type CategoryService struct {
	categoryRepo repositories.CategoryInterface
}

func NewCategoryService(categoryRepo repositories.CategoryInterface) *CategoryService{
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) CreateCategory(category *models.Category) error {
	// NOTE - เช็คว่ามีชื่อ category เป็นค่าว่างไหม
	if category.Name == "" || category.Description == "" {
		return errors.New("Category name and description cannot be empty")
	}

	// NOTE - เช็คว่ามี category ซ้ำไหม
	existingCategory, err := s.categoryRepo.FindByName(category.Name)

	if err != nil {
		return errors.New("Error checking for existing category")
	}

	// NOTE - ถ้ามี category ซ้ำ
	if existingCategory != nil {
		return errors.New("Category already exists")
	}

	// NOTE - สร้าง slug
	category.Slug = slug.Make(category.Name)

	// NOTE - ถ้าไม่มี category ซ้ำก็ทำการสร้าง category ใหม่
	err = s.categoryRepo.Create(category)

	if err != nil {
		return errors.New("Error creating category")
	}
	return nil
}

func (s *CategoryService) UpdateCategory(id uint, category *models.Category) error {
	// NOTE - เช็คว่ามีชื่อ category เป็นค่าว่างไหม
	if category.Name == "" || category.Description == "" {
		return errors.New("Category name and description cannot be empty")
	}

	// NOTE - เช็คว่า category มีอยู่ในระบบไหม
	existingCategory,err  := s.categoryRepo.FindByID(id)

	if err != nil {
		return errors.New("Error finding category")
	}

	if existingCategory == nil {
		return errors.New("Category not found")
	}

	// NOTE - Update Fields ของ category เป็นค่าใหม่
	existingCategory.Slug = slug.Make(category.Name)
	existingCategory.Name = category.Name
	existingCategory.Description = category.Description

	err = s.categoryRepo.Update(existingCategory)

	if err != nil {
		return errors.New("Error updating category")
	}

	return nil
}

func (s *CategoryService) DeleteCategory(id uint) error {
	// NOTE - เช็คว่า category มีอยู่ในระบบไหม
	existingCategory, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("Error finding category")
	}

	if existingCategory == nil {
		return errors.New("Category not found")
	}

	err = s.categoryRepo.Delete(id)
	if err != nil {
		return errors.New("Error deleting category")
	}

	return nil
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, errors.New("Error retrieving categories")
	}

	return categories, nil
}