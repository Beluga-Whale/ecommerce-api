package services

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type CategoryServiceMock struct {
	mock.Mock
}

func NewCategoryServiceMock() *CategoryServiceMock {
	return &CategoryServiceMock{}
}

func (m *CategoryServiceMock) CreateCategory(category *models.Category) error{
	args := m.Called(category)
	return args.Error(0)
}

func (m *CategoryServiceMock) UpdateCategory(id uint, category *models.Category) error{
	args := m.Called(id,category)
	return args.Error(0)
}

func (m *CategoryServiceMock) DeleteCategory(id uint) error{
	args := m.Called(id)
	return args.Error(0)
}

func (m *CategoryServiceMock) GetAllCategories() ([]models.Category, error){
	args := m.Called()

	if categories, ok := args.Get(0).([]models.Category); ok {
		return categories,args.Error(1)
	}

	return nil,args.Error(1)
}
