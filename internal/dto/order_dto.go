package dto

import "github.com/Beluga-Whale/ecommerce-api/internal/models"

type CreateOrderItemDTO struct{
	ProductID uint `json:"productID"`
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
	ProductID       uint    `json:"productID"`
	ProductName     string  `json:"productName"`
	Quantity        uint    `json:"quantity"`
	PriceAtPurchase float64 `json:"priceAtPurchase"`
}
