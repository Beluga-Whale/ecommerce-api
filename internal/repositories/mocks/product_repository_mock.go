package repositories

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func NewProductRepositoryMock() *ProductRepositoryMock {
	return &ProductRepositoryMock{}
}

func (m *ProductRepositoryMock) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) FindByID(id uint) (*models.Product, error){
	args := m.Called(id)

	if product,ok := args.Get(0).(*models.Product); ok {
		return product,nil
	}
	return nil, args.Error(1)
}

func (m *ProductRepositoryMock) FindAll() ([]models.Product, error){
	args := m.Called()

	if products,ok := args.Get(0).([]models.Product); ok {
		return products,nil
	}
	return nil, args.Error(1)
}

func (m *ProductRepositoryMock) Update(product *models.Product) error{
	args := m.Called(product)

	return args.Error(0)
}

func (m *ProductRepositoryMock) Delete(id uint) error{
	args := m.Called(id)

	return args.Error(0)
}

