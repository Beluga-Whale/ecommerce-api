package dto

type ProductCreateDTO struct {
	Name        string              `json:"name" validate:"required,min=2,max=100"`
	Title       string              `json:"title" validate:"required,min=2,max=100"`
	Description string              `json:"description" validate:"required,min=10`
	Images      []ProductImageDTO   `json:"images" validate:"required,dive"`
	Variants    []ProductVariantDTO `json:"variants" validate:"required,dive"`
	IsFeatured  bool                `json:"isFeatured" validate:"omitempty"`
	IsOnSale    bool                `json:"isOnSale" validate:"omitempty"`
	SalePrice   *float64            `json:"salePrice"`
	CategoryID  uint                `json:"categoryID"`
}

type ProductCreateResponseDTO struct {
	ID          uint                `json:"id" validate:"required"`
	Name        string              `json:"name" validate:"required,min=2,max=100"`
	Title       string              `json:"title" validate:"required,min=2,max=100"`
	Description string              `json:"description" validate:"required,min=10`
	Images      []ProductImageDTO   `json:"images" validate:"required,dive"`
	Variants    []ProductVariantDTO `json:"variants" validate:"required,dive"`
	IsFeatured  bool                `json:"isFeatured" validate:"omitempty"`
	IsOnSale    bool                `json:"isOnSale" validate:"omitempty"`
	SalePrice   *float64            `json:"salePrice"`
	CategoryID  uint                `json:"categoryID"`
}

type ProductUpdateDTO struct {
	Name        string              `json:"name" validate:"required,min=2,max=100"`
	Title       string              `json:"title" validate:"required,min=2,max=100"`
	Description string              `json:"description" validate:"required,min=10`
	Images      []ProductImageDTO   `json:"images" validate:"required,dive"`
	Variants    []ProductVariantDTO `json:"variants" validate:"required,dive"`
	IsFeatured  bool                `json:"isFeatured" validate:"omitempty"`
	IsOnSale    bool                `json:"isOnSale" validate:"omitempty"`
	SalePrice   *float64            `json:"salePrice"`
	CategoryID  uint                `json:"categoryID"`
}

type ProductUpdateResponseDTO struct {
	Name        string              `json:"name" validate:"required,min=2,max=100"`
	Title       string              `json:"title" validate:"required,min=2,max=100"`
	Description string              `json:"description" validate:"required,min=10`
	Images      []ProductImageDTO   `json:"images" validate:"required,dive"`
	Variants    []ProductVariantDTO `json:"variants" validate:"required,dive"`
	IsFeatured  bool                `json:"isFeatured" validate:"omitempty"`
	IsOnSale    bool                `json:"isOnSale" validate:"omitempty"`
	SalePrice   *float64            `json:"salePrice"`
	CategoryID  uint                `json:"categoryID"`
}

type ProductVariantDTO struct {
	Size       string  `json:"size" validate:"required"`
	Stock      int     `json:"stock" validate:"required,min=0"`
	SKU        string  `json:"sku" validate:"required"`
	Price      float64 `json:"price" validate:"required,gt=0"`
	FinalPrice float64 `json:"finalPrice"`
}

type ProductImageDTO struct {
	URL string `json:"url" validate:"required"`
}