package models

import "gorm.io/gorm"

type ProductVariant struct {
	gorm.Model
	ProductID uint	
	Size      string
	Stock     int
	SKU       string
	Price     float64
}