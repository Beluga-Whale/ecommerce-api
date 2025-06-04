package services

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/repositories"
	"github.com/Beluga-Whale/ecommerce-api/internal/utils"
)

type OrderServiceInterface interface {
}

type OrderService struct {
	orderRepo repositories.OrderRepositoryInterface
	productUtil utils.ProductInterface
}

func NewOrderService(orderRepo repositories.OrderRepositoryInterface,productUtil utils.ProductInterface) *OrderService {
	return &OrderService{orderRepo: orderRepo,productUtil:productUtil }
}

