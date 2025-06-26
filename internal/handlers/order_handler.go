package handlers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

type OrderHandlerInterface interface{
	CreateOrder(c *fiber.Ctx) error
	UpdateStatusOrder(c *fiber.Ctx) error
	GetOrderByID(c *fiber.Ctx) error 
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
		return JSONError(c, fiber.StatusBadRequest,"Invalid request body")
	}

	// NOTE - เอา UserIDจาก local
	// NOTE - ดึง userID จาก Locals แล้วแปลง string -> uint
	userIDStr, ok := c.Locals("userID").(string)

	if !ok {
		return JSONError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return JSONError(c, fiber.StatusInternalServerError, "Invalid user ID format")
	}

	order,err := h.OrderService.CreateOrder(uint(userIDUint), req)

	if err != nil{
		return JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	response := dto.OrderResponseDTO{
		OrderID: order.ID,
		Status:  order.Status,
		TotalPrice: order.TotalPrice,
		User : order.User.ID,
		FullName: order.FullName,
		Phone:    order.Phone,
		Address:  order.Address,
		Province: order.Province,
		District: order.District,
		Subdistrict: order.Subdistrict,
		Zipcode: order.Zipcode,
		Coupon: order.Coupon.ID,
		Items: []dto.OrderItemResponseDTO{},
	}

	for _, item := range order.OrderItem {
		response.Items = append(response.Items, dto.OrderItemResponseDTO{
			VariantID:       item.ProductVariantID,
			ProductName:     item.ProductVariant.Product.Name, 
			Quantity:        item.Quantity,
			Size: 			 item.ProductVariant.Size,	
			PriceAtPurchase: item.PriceAtPurchase,
		})
	}


	return JSONSuccess(c, fiber.StatusCreated, "Order created successfully", response)
	
}

func (h *OrderHandler) UpdateStatusOrder(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	expectedToken := "Bearer " + os.Getenv("STRIPE_WEBHOOK_SECRET")

	fmt.Println("Auth Header:", authHeader)

	if authHeader != expectedToken {
		return JSONError(c, fiber.StatusUnauthorized, "Unauthorized - missing token")	}
	var req dto.UpdateStatusOrderDTO

	if err := c.BodyParser(&req); err != nil {
		return JSONError(c, fiber.StatusBadRequest,"Invalid request body")
	}

	if err := h.OrderService.UpdateStatusOrder(&req.OrderId,req.Status,req.UserId); err !=nil{
		return JSONError(c, fiber.StatusInternalServerError, err.Error())
	}

	return JSONSuccess(c,fiber.StatusOK,"Update Status Order Sucess",nil)

}

func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	userId := c.QueryInt("userId",0)
	orderId, err := c.ParamsInt("id")
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid order ID")
	}

	order,err := h.OrderService.GetOrderByID(uint(orderId),uint(userId))
	if err != nil {
		return JSONError(c, fiber.StatusInternalServerError,  "Error go get order by ID")
	}

	if order == nil {
		return JSONError(c, fiber.StatusNotFound, "Order not found")
	}

	var orderItems []dto.OrderItemResponseDTO

	for _, item := range order.OrderItem {
		orderItems = append(orderItems, dto.OrderItemResponseDTO{
			VariantID:       item.ProductVariantID,
			ProductName:     item.ProductVariant.Product.Name,
			Size:            item.ProductVariant.Size,
			Quantity:        item.Quantity,
			PriceAtPurchase: item.PriceAtPurchase,
		})
	}

	return JSONSuccess(c,fiber.StatusOK,"Get Order By ID Success", dto.OrderByIDResponseDTO{
		OrderID:   order.ID,
		Status:    order.Status,
		FullName:  order.FullName,
		Phone:     order.Phone,
		Address:   order.Address,
		Province:  order.Province,
		District:  order.District,
		Subdistrict: order.Subdistrict,
		Zipcode:   order.Zipcode,
		Coupon:     order.Coupon,
		OrderItem: orderItems,
		PaymentExpireAt: order.PaymentExpireAt.Format("2006-01-02 15:04:05"),
	})
}

