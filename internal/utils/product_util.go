package utils

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
)

type ProductInterface interface {
	FindProductID(products []models.Product, id uint) *models.Product
}

type Product_Util struct{}

func NewProductUtil() *Product_Util {
	return &Product_Util{}
}

func (h *Product_Util) FindProductID(products []models.Product, id uint) *models.Product {
	for _, p := range products {
		if p.ID == id {
			return &p
		}
	}
	return nil
}
