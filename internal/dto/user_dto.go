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

// NOTE - Category DTOs

// NOTE - Create, Update, Delete Category DTOs

type CategoryCreateDTO struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type CategoryCreateResponseDTO struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	Slug string `json:"slug"`
}

type UpdateCategoryDTO struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type UpdateCategoryResponseDTO struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	Slug string `json:"slug"`
}