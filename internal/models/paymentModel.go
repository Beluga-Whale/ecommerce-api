package models

import (
	"time"

	"gorm.io/gorm"
)

type StatusPayment string

const (
	Payed StatusPayment = "payed"
	Failed StatusPayment = "failed"
)


type Payment struct {
	gorm.Model
	OrderID uint //NOTE -
	Order Order `gorm:"foreignKey:OrderID"`
	Amount float64
	Status StatusPayment `gorm:"type:payment_status"`
	PaidAt time.Time
}