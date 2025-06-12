package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	UserRole Role = "user"
	AdminRole Role = "admin"
)

type User struct {
	gorm.Model
	FirstName string
	LastName string
	Email string
  	Password string
	Phone string
	Date time.Time
  	Role Role `gorm:"type:role;default:'user'"`
}