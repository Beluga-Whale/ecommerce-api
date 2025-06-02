package services

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type ProductServiceMock struct {
	mock.Mock
}

func NewProductServiceMock() *ProductServiceMock {
	return &ProductServiceMock{}
}

func (m *ProductServiceMock) CreateProduct(product *models.Product) error {
	args := m.Called(product)
	return  args.Error(0)
}

func (m *ProductServiceMock) UpdateProduct(id uint, product *models.Product) error {
	args := m.Called(id, product)
	return  args.Error(0)
}

func (m *ProductServiceMock) DeleteProduct(id uint) error {
	args := m.Called(id)
	return  args.Error(0)
}

func (m *ProductServiceMock) GetProductByID(id uint) (*models.Product, error) {
	args := m.Called(id)
	if product,ok := args.Get(0).(*models.Product) ; ok {
		return product,nil
	}
	return  nil,args.Error(1)
}

func (m *ProductServiceMock) GetAllProducts() ([]models.Product, error) {
	args := m.Called()
	if products,ok := args.Get(0).([]models.Product) ; ok {
		return products,nil
	}
	return  nil,args.Error(1)
}