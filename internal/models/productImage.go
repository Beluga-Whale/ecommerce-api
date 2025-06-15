package models

import "gorm.io/gorm"

type ProductImage struct {
	gorm.Model
	URL string
	ProductID uint
	Product Product
}