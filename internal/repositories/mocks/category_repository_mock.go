package repositories

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type CategoryRepositoryMock struct {
	mock.Mock
}

func NewCategoryRepositoryMock() *CategoryRepositoryMock {
	return &CategoryRepositoryMock{}
}

func (m *CategoryRepositoryMock) Create(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) Update(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) FindAll() ([]models.Category,error) {
	args := m.Called()
	if categories, ok := args.Get(0).([]models.Category); ok {
		return categories, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryRepositoryMock) FindByName(name string) (*models.Category,error) {
	args := m.Called(name)
	if category, ok := args.Get(0).(*models.Category); ok {
		return category, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CategoryRepositoryMock) FindByID(id uint) (*models.Category,error) {
	args := m.Called(id)
	if category, ok := args.Get(0).(*models.Category); ok {
		return category, args.Error(1)
	}
	return nil, args.Error(1)
}