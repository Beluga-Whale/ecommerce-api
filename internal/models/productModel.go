package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name string
	Description string `gorm:"type:text"`
	Image string
	IsFeatured bool
	IsOnSale bool
	SalePrice *float64
	CategoryID uint //NOTE FK
	Category Category `gorm:"foreignKey:CategoryID"`
	Variants []ProductVariant `gorm:"foreignKey:ProductID"`
}