package dto

import "github.com/Beluga-Whale/ecommerce-api/internal/models"

type CreateOrderItemDTO struct{
	VariantID uint `json:"variantID"`
	Quantity uint `json:"quantity"`
}

type CreateOrderRequestDTO struct {
	Counpon *uint `json:"counponID"`
	Note string `json:"note"`
	Phone string `json:phone`
	Address string `json:string`
	Items []CreateOrderItemDTO `json:"items"`
}

type OrderResponseDTO struct {
	OrderID    uint                    `json:"orderID"`
	Status     models.Status           `json:"status"`
	TotalPrice float64                 `json:"totalPrice"`
	Items      []OrderItemResponseDTO  `json:"items"`
}

type OrderItemResponseDTO struct {
	VariantID       uint    `json:"variantID"`
	ProductName     string  `json:"productName"`
	Size            string  `json:"size"`
	Quantity        uint    `json:"quantity"`
	PriceAtPurchase float64 `json:"priceAtPurchase"`
}
