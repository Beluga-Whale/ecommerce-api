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

func (m *ProductRepositoryMock) FindAll(page uint, limit uint, minPrice int64, maxPrice int64, searchName string, category string) ([]models.Product, int64, error) {
	args := m.Called(page, limit, minPrice, maxPrice, searchName, category)

	var products []models.Product
	if res, ok := args.Get(0).([]models.Product); ok {
		products = res
	}

	var pageTotal int64
	if pt, ok := args.Get(1).(int64); ok {
		pageTotal = pt
	}

	err := args.Error(2)

	return products, pageTotal, err
}

func (m *ProductRepositoryMock) Update(product *models.Product) error{
	args := m.Called(product)

	return args.Error(0)
}

func (m *ProductRepositoryMock) Delete(id uint) error{
	args := m.Called(id)

	return args.Error(0)
}

