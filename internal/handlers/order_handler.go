package handlers

import (
	"strconv"

	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
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
	// NOTE - สร้างตัวแปรเก็บ reqBody
	var req dto.CreateOrderRequestDTO

	if err := c.BodyParser(&req); err !=nil {
		return JSONError(c, fiber.StatusBadRequest,"Invalid request bod")
	}

	// NOTE - เอา UserIDจาก local
	// ดึง userID จาก Locals แล้วแปลง string -> uint
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return JSONError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return JSONError(c, fiber.StatusInternalServerError, "Invalid user ID format")
	}

	order,err := h.OrderService.CreateOrder(uint(userIDUint), req)

	response := dto.OrderResponseDTO{
		OrderID: order.ID,
		Status:  order.Status,
		TotalPrice: order.TotalPrice,
		Items: []dto.OrderItemResponseDTO{},
	}

	for _, item := range order.OrderItem {
		response.Items = append(response.Items, dto.OrderItemResponseDTO{
			ProductID:       item.ProductID,
			ProductName:     item.Product.Name, // ต้อง preload มาก่อน
			Quantity:        item.Quantity,
			PriceAtPurchase: item.PriceAtPurchase,
		})
	}

	if err != nil{
		return JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return JSONSuccess(c, fiber.StatusCreated, "Order created successfully", response)
	
}