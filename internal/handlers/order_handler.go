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
	UpdateOrderStatusByUser(c *fiber.Ctx) error
	GetAllOrders(c *fiber.Ctx) error
	UpdateOrderStatusByAdmin(c *fiber.Ctx) error
	GetSummary(c *fiber.Ctx) error
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
		Zipcode:    order.Zipcode,
		Coupon:     order.Coupon,
		OrderItem:  orderItems,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt.Format("2006-01-02 15:04:05"),
		PaymentExpireAt: order.PaymentExpireAt.Format("2006-01-02 15:04:05"),
	})
}

func (h*OrderHandler) GetAllOrderByUserId(c *fiber.Ctx) error {

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

	orderAll, err :=h.OrderService.GetAllOrderByUserId(uint(userIDUint))

	var orderList []dto.OrderListResponseDTO

	for _, order := range orderAll {
	orderList = append(orderList, dto.OrderListResponseDTO{
		OrderID:    order.ID,
		TotalPrice: order.TotalPrice,
		Status:     string(order.Status),
		ItemCount:  len(order.OrderItem),
		CreatedAt:  order.CreatedAt.Format("2006-01-02 15:04:05"),
	})
	}

	return JSONSuccess(c, fiber.StatusOK, "Get All Orders Success", orderList)
}

func (h *OrderHandler) UpdateOrderStatusByUser(c *fiber.Ctx) error {

	orderID, err := c.ParamsInt("id")
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	var req dto.UpdateStatusByUserOrderDTO

	if err := c.BodyParser(&req); err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return JSONError(c, fiber.StatusUnauthorized, "Unauthorized")
	}
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	orderIDUint := uint(orderID)
	err = h.OrderService.UpdateStatusByUser(uint(userIDUint), &orderIDUint, req.Status)
	
	if err != nil {
    	return JSONError(c, fiber.StatusInternalServerError, fmt.Sprintf("Update status failed: %v", err))
	}

	return JSONSuccess(c, fiber.StatusOK, "Update status success", nil)
}

func (h *OrderHandler)	GetAllOrders(c *fiber.Ctx) error {
	orders, err := h.OrderService.GetAllOrdersAdmin()

	if err != nil {
		return JSONError(c, fiber.StatusInternalServerError,"Get Order Error")	
	}

	var orderResponse []dto.OrderListDataTableDTOResponse

	for _, item := range orders {
		var orderItemResponse []dto.OrderItemResponseDTO

		for _, v := range item.OrderItem {
			orderItemResponse = append(orderItemResponse, dto.OrderItemResponseDTO{
				VariantID:       v.ProductVariantID,
				ProductName:     v.ProductVariant.Product.Name,
				Size:            v.ProductVariant.Size,
				Quantity:        v.Quantity,
				PriceAtPurchase: v.PriceAtPurchase,
			})
		}

		orderResponse = append(orderResponse, dto.OrderListDataTableDTOResponse{
			OrderID:    item.ID,
			CreatedAt:  item.CreatedAt.Format("2006-01-02 15:04:05"),
			UserName:   item.FullName,
			Status:     item.Status,
			TotalPrice: item.TotalPrice,
			OrderItem:  orderItemResponse,
		})
	}

	return JSONSuccess(c,fiber.StatusOK,"Get All Order Success",orderResponse)
}

func (h *OrderHandler) UpdateOrderStatusByAdmin(c *fiber.Ctx) error{
	
	orderID, err := c.ParamsInt("id")
	if err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	var req dto.UpdateStatusByUserOrderDTO
	if err := c.BodyParser(&req); err != nil {
		return JSONError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	orderIDUint := uint(orderID)
	err = h.OrderService.UpdateStatusByAdmin( &orderIDUint, req.Status)
	
	if err != nil {
    	return JSONError(c, fiber.StatusInternalServerError, fmt.Sprintf("Update status failed: %v", err))
	}

	return JSONSuccess(c, fiber.StatusOK, "Update status success", nil)

}

func (h *OrderHandler) GetSummary(c *fiber.Ctx) error {
	summary, err := h.OrderService.GetDashboardSummary()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch summary",
		})
	}

	return JSONSuccess(c, fiber.StatusOK, "Get summary success", dto.DashboardSummaryDTO{
		OrderTotal: summary.OrderTotal,
		OrdersThisMonth:summary.OrdersThisMonth,
		OrdersLastMonth:summary.OrdersLastMonth,
		OrderGrowthPercent:summary.OrderGrowthPercent,
		RevenueThisMonth:summary.RevenueThisMonth,
		RevenueLastMonth:summary.RevenueLastMonth,
		RevenueGrowthPercent:summary.RevenueGrowthPercent,
		CustomersThisMonth:summary.CustomersThisMonth,
		CustomersLastMonth:summary.CustomersLastMonth,
		CustomerGrowthPercent:summary.CustomerGrowthPercent,
	})
}
