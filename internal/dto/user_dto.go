package dto

// NOTE - User DTOs
type RegisterRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
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
}
