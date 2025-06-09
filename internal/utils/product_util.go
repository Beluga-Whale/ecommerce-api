package utils

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
)

type ProductInterface interface {
	FindProductVariantID(products []models.ProductVariant, productID uint) *models.ProductVariant
}

type Product_Util struct{}

func NewProductUtil() *Product_Util {
	return &Product_Util{}
}

func (h *Product_Util) FindProductVariantID(products []models.ProductVariant, productID uint) *models.ProductVariant {
	for _, p := range products {
		if p.ID == productID {
			return &p
		}
	}
	return nil
}
