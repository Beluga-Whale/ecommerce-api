package models

import "gorm.io/gorm"

type Status string

const (
	Pending Status = "pending"
	Paid Status = "paid"
	Shipped Status = "shipped"
	Cancel Status = "cancel"
)


type Order struct {
	gorm.Model
	UserID uint //NOTE FK
	User User `gorm:"foreignKey:UserID"`
	CouponID *uint //NOTE FK แล้วสามารถเป็นค่าว่างได้ nullable
	Coupon Coupon `gorm:"foreignKey:CouponID"`
	Status Status `gorm:"type:status;default:'pending'"`
	Note string
	Phone string
	Address string
	TotalPrice float64
	OrderItem []OrderItem `gorm:"foreignKey:OrderID"`

}