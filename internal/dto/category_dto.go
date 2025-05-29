package dto

// NOTE - Category DTOs

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