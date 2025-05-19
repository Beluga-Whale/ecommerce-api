package models

import "gorm.io/gorm"

type Role string

const (
	UserRole Role = "user"
	AdminRole Role = "admin"
)

type User struct {
	gorm.Model
	Name string
	Email string
  	Password string
  	Role Role `gorm:"type:role;default:'user'"`
}