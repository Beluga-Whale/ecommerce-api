package dto

type CreatePaymentIntentRequestDTO struct {
	Amount  int64 `json:"amount"`
	OrderID uint  `json:"orderId"`
	UserID  uint  `json:"userId"`
}