package utils

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type ProductUtilMock struct {
	mock.Mock
}


func NewProductUtilMock() *ProductUtilMock {
	return &ProductUtilMock{}
}

func (m *ProductUtilMock) FindProductVariantID(products []models.ProductVariant, productID uint) *models.ProductVariant {
	args := m.Called(products,productID)

	return args.Get(0).(*models.ProductVariant)
}