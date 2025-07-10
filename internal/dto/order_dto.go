package dto

import "github.com/Beluga-Whale/ecommerce-api/internal/models"

type CreateOrderItemDTO struct{
	VariantID uint `json:"variantID"`
	Quantity uint `json:"quantity"`
}

type CreateOrderRequestDTO struct {
	Counpon *uint `json:"counponID"`
	FullName string `json:"fullName"`
	Phone string `json:phone`
	Address string `json:address`
	Province string `json:province`
	District string `json:district`
	Subdistrict string `json:subdistrict`
	Zipcode string `json:zipcode`
	Items []CreateOrderItemDTO `json:"items"`
}

type OrderResponseDTO struct {
	OrderID    uint                    `json:"orderID"`
	User       uint                    `json:"user"`
	FullName   string  				   `json:"fullName"`
	Phone 	   string 				   `json:phone`
	Address    string 				   `json:address`
	Province   string 				   `json:province`
	District   string 				   `json:district`
	Subdistrict string 				   `json:subdistrict`
	Zipcode    string 				   `json:zipcode`
	Coupon     uint		   			   `json:"coupon"`
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
	ProductID       uint    `json:"productId"`
}


type UpdateStatusOrderDTO struct {
	OrderId uint           `json:"orderId"`
	Status  models.Status  `json:"status"`
	UserId  uint		   `json:"userId"`
}

type OrderByIDResponseDTO struct {
	OrderID    uint                    `json:"orderID"`
	Status    models.Status            `json:"status"`
	FullName   string  				   `json:"fullName"`
	Phone 	   string 				   `json:"phone"`
	Address    string 				   `json:"address"`
	Province   string 				   `json:"province"`
	District   string 				   `json:"district"`
	Subdistrict string 				   `json:"subdistrict"`
	Zipcode    string 				   `json:"zipcode"`
	Coupon     models.Coupon           `json:"coupon"`
	OrderItem  []OrderItemResponseDTO  `json:"orderItem"`
	TotalPrice  float64 `json:"totalPrice"`
	PaymentExpireAt string             `json:"paymentExpireAt"`
	CreatedAt string `json:"createdAt"`
}

type OrderListResponseDTO struct {
	OrderID     uint    `json:"orderID"`
	TotalPrice  float64 `json:"totalPrice"`
	Status      string  `json:"status"`
	ItemCount   int     `json:"itemCount"`
	CreatedAt   string  `json:"createdAt"`
}

type UpdateStatusByUserOrderDTO struct {
	Status  models.Status  `json:"status" validate:"required"`
}

type OrderListDataTableDTOResponse struct {
	OrderID     uint    `json:"orderID"`
	CreatedAt  string  `json:"createdAt"`
	UserName   string  `json:"userName"`
	Status     models.Status `json:"status"`
	TotalPrice  float64 `json:"totalPrice"`
	OrderItem  []OrderItemResponseDTO  `json:"orderItem"`
}