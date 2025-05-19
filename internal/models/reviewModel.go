package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	UserID uint //NOTE - FK
	User User `gorm:"foreignKey:UserID"`
	ProductID uint //NOTE - FK
	Product Product `gorm:"foreignKey:ProductID"`
	Rating uint 
	Comment string
}