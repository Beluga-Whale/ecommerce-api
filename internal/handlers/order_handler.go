package handlers

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

type OrderHandlerInterface interface{
	CreateOrder(c *fiber.Ctx) error
}

type OrderHandler struct {
	OrderService services.OrderServiceInterface
}

func NewOrderHandler(OrderService services.OrderServiceInterface) *OrderHandler {
	return &OrderHandler{OrderService:OrderService}
}


func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	return nil
}