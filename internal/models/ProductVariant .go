package models

import "gorm.io/gorm"

type ProductVariant struct {
	gorm.Model
	ProductID uint	
    Product   Product
	Size      string
	Stock     int
	SKU       string
	Price     float64
}