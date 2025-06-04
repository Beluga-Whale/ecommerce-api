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

func (m *ProductServiceMock) GetAllProducts( page uint, limit uint) ([]models.Product, int64, error)  {
	args := m.Called(page, limit)

	var products []models.Product
	if res,ok := args.Get(0).([]models.Product);ok {
		products = res
	}

	var pageTotal int64
	if pt, ok := args.Get(1).(int64); ok {
		pageTotal = pt
	}

	return  products, pageTotal,args.Error(2)
}