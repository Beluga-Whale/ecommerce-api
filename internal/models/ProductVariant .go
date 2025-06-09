package models

import "gorm.io/gorm"

type ProductVariant struct {
	gorm.Model
	ProductID uint	
	Product Product `gorm:"foreignKey:ProductID"`
	Size      string
	Stock     int
	SKU       string
}