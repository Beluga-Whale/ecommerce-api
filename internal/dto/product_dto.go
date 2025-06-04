package dto

type ProductCreateDTO struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Description string  `json:"description" validate:"required,min=10,max=500"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Image       string  `json:"image" validate:"required,url"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	IsFeatured  bool    `json:"isFeatured" validate:"omitempty"`
	IsOnSale    bool    `json:"isOnSale" validate:"omitempty"`
	SalePrice   *float64 `json:"salePrice"`
	CategoryID  uint    `json:"categoryID"`
}

type ProductCreateResponseDTO struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Description string  `json:"description" validate:"required,min=10,max=500"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Image       string  `json:"image" validate:"required,url"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	IsFeatured  bool    `json:"isFeatured" validate:"omitempty"`
	IsOnSale    bool    `json:"isOnSale" validate:"omitempty"`
	SalePrice   *float64 `json:"salePrice"`
	CategoryID  uint    `json:"categoryID"`
}

type ProductUpdateDTO struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Description string  `json:"description" validate:"required,min=10,max=500"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Image       string  `json:"image" validate:"required,url"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	IsFeatured  bool    `json:"isFeatured" validate:"omitempty"`
	IsOnSale    bool    `json:"isOnSale" validate:"omitempty"`
	SalePrice   *float64 `json:"salePrice"`
	CategoryID  uint    `json:"categoryID"`
}

type ProductUpdateResponseDTO struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Description string  `json:"description" validate:"required,min=10,max=500"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Image       string  `json:"image" validate:"required,url"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	IsFeatured  bool    `json:"isFeatured" validate:"omitempty"`
	IsOnSale    bool    `json:"isOnSale" validate:"omitempty"`
	SalePrice   *float64 `json:"salePrice"`
	CategoryID  uint    `json:"categoryID"`
}
