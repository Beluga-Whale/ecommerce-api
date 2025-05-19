package models

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model
	Code string
	DiscountAmount float64
	ExpiredAt time.Time
}