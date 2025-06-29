package dto

import "time"

// NOTE - User DTOs
type RegisterRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FirstName string `json:"firstName" validate:"required,min=2,max=100"`
	LastName string `json:"lastName" validate:"required,min=2,max=100"`
	Phone string `json:"phone" validate:"required,min=10,max=10"`
	BirthDate time.Time `json:"birthDate"`
}

type RegisterResponseDTO struct {
	Message string `json:"message"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponseDTO struct {
	Token string `json:"token"`
	UserID uint   `json:"userId"`
}

type UserProfileDTO struct {
	UserID uint `json:"userId"`
	Email    string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required,min=2,max=100"`
	LastName string `json:"lastName" validate:"required,min=2,max=100"`
	Phone string `json:"phone" validate:"required,min=10,max=10"`
	BirthDate time.Time `json:"birthDate"`
	Avatar string `json:"avatar"`
}


type UserUpdateProfileDTO struct {
	FirstName   *string    `json:"firstName"`
	LastName    *string    `json:"lastName"`
	Phone       *string    `json:"phone"`
	BirthDate   *time.Time `json:"birthDate"`
	Avatar      *string `json:"avatar"`
}