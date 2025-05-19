package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	OrderID uint //NOTE - FK
	Order Order `gorm:"foreignKey:OrderID"`
	ProductID uint //NOTE - FK
	Product Product `gorm:"foreignKey:ProductID"`
	Quantity uint 
	PriceAtPurchase float64
}