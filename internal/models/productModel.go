package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name string
	Description string
	Price string
	Image string
	Stock int 
	IsFeatured bool
	IsOnSale bool
	CategoryID uint //NOTE FK
	Category Category `gorm:"foreignKey:CategoryID"`
}