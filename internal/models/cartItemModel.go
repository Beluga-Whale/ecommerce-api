package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	UserID uint //NOTE FK
	User User `gorm:"foreignKey:UserID"`
	ProductID uint //NOTE FK
	Product Product `gorm:"foreignKey:ProductID"`
	Quantity uint
}

