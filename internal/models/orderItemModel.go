package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	OrderID uint //NOTE - FK
	Order Order `gorm:"foreignKey:OrderID"`
	ProductVariantID uint //NOTE - FK
	ProductVariant   ProductVariant `gorm:"foreignKey:ProductVariantID"`
	Quantity uint 
	PriceAtPurchase float64
}