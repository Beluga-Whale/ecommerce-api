package handlers_test

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Beluga-Whale/ecommerce-api/internal/handlers"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	services "github.com/Beluga-Whale/ecommerce-api/internal/services/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateOrder(t *testing.T) {
	t.Run("CreateOrder",func(t *testing.T) {
		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			Status: "pending",
		}

		orderService := services.NewOrderServiceMock()

		orderHandler := handlers.NewOrderHandler(orderService)

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		orderService.On("CreateOrder",mock.Anything,mock.Anything).Return(&orderMock,nil)


		
		app := fiber.New()
		app.Post("/user/order",testMiddleware,orderHandler.CreateOrder)

		reqBody:= []byte(`{
			"fullName": "T-Shirt Update",
			"phone": "0987678976",
			"address": "1",
			"province": "1",
			"district": "1",
			"subdistrict": "1",
			"zipcode": "1",
			"items": [
				{
					"variantId": 92,
					"quantity": 1
				}
				
			]
		}`)

		req :=httptest.NewRequest("POST","/user/order",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Order created successfully")
	})

	t.Run("Invalid request body",func(t *testing.T) {
		orderService := services.NewOrderServiceMock()

		orderHandler := handlers.NewOrderHandler(orderService)

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}
		
		app := fiber.New()
		app.Post("/user/order",testMiddleware,orderHandler.CreateOrder)

		req :=httptest.NewRequest("POST","/user/order",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request body")
	})

	t.Run("Unauthorized",func(t *testing.T) {
		orderService := services.NewOrderServiceMock()

		orderHandler := handlers.NewOrderHandler(orderService)

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", nil)
			return c.Next()
		}

		app := fiber.New()
		app.Post("/user/order",testMiddleware,orderHandler.CreateOrder)

		reqBody:= []byte(`{
			"fullName": "T-Shirt Update",
			"phone": "0987678976",
			"address": "1",
			"province": "1",
			"district": "1",
			"subdistrict": "1",
			"zipcode": "1",
			"items": [
				{
					"variantId": 92,
					"quantity": 1
				}
				
			]
		}`)

		req :=httptest.NewRequest("POST","/user/order",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Unauthorized")
	})

	t.Run("Unauthorized InvalidUserIDFormat",func(t *testing.T) {
		orderService := services.NewOrderServiceMock()

		orderHandler := handlers.NewOrderHandler(orderService)

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "faileUserID")
			return c.Next()
		}

		app := fiber.New()
		app.Post("/user/order",testMiddleware,orderHandler.CreateOrder)

		reqBody:= []byte(`{
			"fullName": "T-Shirt Update",
			"phone": "0987678976",
			"address": "1",
			"province": "1",
			"district": "1",
			"subdistrict": "1",
			"zipcode": "1",
			"items": [
				{
					"variantId": 92,
					"quantity": 1
				}
				
			]
		}`)

		req :=httptest.NewRequest("POST","/user/order",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid user ID format")
	})

	t.Run("Error to create order",func(t *testing.T) {

		orderService := services.NewOrderServiceMock()

		orderHandler := handlers.NewOrderHandler(orderService)

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		orderService.On("CreateOrder",mock.Anything,mock.Anything).Return(nil,errors.New("Error to create order"))


		
		app := fiber.New()
		app.Post("/user/order",testMiddleware,orderHandler.CreateOrder)

		reqBody:= []byte(`{
			"fullName": "T-Shirt Update",
			"phone": "0987678976",
			"address": "1",
			"province": "1",
			"district": "1",
			"subdistrict": "1",
			"zipcode": "1",
			"items": [
				{
					"variantId": 92,
					"quantity": 1
				}
				
			]
		}`)

		req :=httptest.NewRequest("POST","/user/order",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Error to create order")
	})

}

func TestUpdateStatusOrder(t *testing.T) {
	t.Run("UpdateStatusOrder success", func(t *testing.T) {
		validToken := "Bearer my-secret"
		_ = os.Setenv("STRIPE_WEBHOOK_SECRET", "my-secret")

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)

		app := fiber.New()
		app.Post("/user/order", orderHandler.UpdateStatusOrder)

		validBody := []byte(`{
			"orderId": 1,
			"status": "paid",
			"userId": 99
		}`)
		orderService.On("UpdateStatusOrder", mock.Anything, models.Status("paid") , uint(99)).Return(nil)

		req := httptest.NewRequest("POST", "/user/order", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", validToken)

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Update Status Order Sucess")
	})

	t.Run("Unauthorized - missing token", func(t *testing.T) {
		_ = os.Setenv("STRIPE_WEBHOOK_SECRET", "my-secret")

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)

		app := fiber.New()
		app.Post("/user/order", orderHandler.UpdateStatusOrder)

		validBody := []byte(`{
			"orderId": 1,
			"status": "paid",
			"userId": 99
		}`)
		req := httptest.NewRequest("POST", "/user/order", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Unauthorized - missing token")
	})

	t.Run("Invalid request body", func(t *testing.T) {
		validToken := "Bearer my-secret"
		_ = os.Setenv("STRIPE_WEBHOOK_SECRET", "my-secret")

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)

		app := fiber.New()
		app.Post("/user/order", orderHandler.UpdateStatusOrder)
		req := httptest.NewRequest("POST", "/user/order", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", validToken)

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request body")
	})

	t.Run("Error UpdateStatusOrder", func(t *testing.T) {
		validToken := "Bearer my-secret"
		_ = os.Setenv("STRIPE_WEBHOOK_SECRET", "my-secret")

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)

		app := fiber.New()
		app.Post("/user/order", orderHandler.UpdateStatusOrder)

		validBody := []byte(`{
			"orderId": 1,
			"status": "paid",
			"userId": 99
		}`)

		orderService.On("UpdateStatusOrder",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(errors.New("Error to update status order"))


		req := httptest.NewRequest("POST", "/user/order", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", validToken)

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Error to update status order")
	})
}

func TestGetOrderByID(t *testing.T) {
	t.Run("GetOrder Success",func(t *testing.T) {
		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			Status: "pending",
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		orderService.On("GetOrderByID",mock.Anything,mock.Anything).Return(&orderMock,nil)

		app := fiber.New()
		app.Get("/user/order/:id", orderHandler.GetOrderByID)

		req := httptest.NewRequest("GET", "/user/order/1",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Get Order By ID Success")
	})

	t.Run("Invalid order ID",func(t *testing.T) {
		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			Status: "pending",
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		orderService.On("GetOrderByID",mock.Anything,mock.Anything).Return(&orderMock,nil)

		app := fiber.New()
		app.Get("/user/order/:id", orderHandler.GetOrderByID)

		req := httptest.NewRequest("GET", "/user/order/failID",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid order ID")
	})

	t.Run("Error GetOrderByID",func(t *testing.T) {
		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		orderService.On("GetOrderByID",mock.Anything,mock.Anything).Return(nil,errors.New("Error go get order by ID"))

		app := fiber.New()
		app.Get("/user/order/:id", orderHandler.GetOrderByID)

		req := httptest.NewRequest("GET", "/user/order/1",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Error go get order by ID")
	})

	t.Run("GetOrder Success",func(t *testing.T) {
		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		orderService.On("GetOrderByID",mock.Anything,mock.Anything).Return(nil,nil)

		app := fiber.New()
		app.Get("/user/order/:id", orderHandler.GetOrderByID)

		req := httptest.NewRequest("GET", "/user/order/1",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Order not found")
	})
}

func TestUpdateOrderStatusByUser(t *testing.T) {
	t.Run("UpdateOrderStatusByUser Success",func(t *testing.T) {
		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			Status: "pending",
		}

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		orderService.On("FindOrderById",mock.Anything).Return(&orderMock,nil)
		orderService.On("UpdateStatusByUser",mock.Anything,mock.Anything,mock.Anything).Return(nil)

		app := fiber.New()
		app.Patch("/user/order/:id/status", testMiddleware,orderHandler.UpdateOrderStatusByUser)

		reqBody:= []byte(`{
			"status":"pending"
		}`)

		req := httptest.NewRequest("PATCH", "/user/order/1/status",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Update status success")
	})

	t.Run("Error invalid orderID",func(t *testing.T) {
		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)

		app := fiber.New()
		app.Patch("/user/order/:id/status", testMiddleware,orderHandler.UpdateOrderStatusByUser)

		reqBody:= []byte(`{
			"status":"pending"
		}`)

		req := httptest.NewRequest("PATCH", "/user/order/failID/status",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid product ID")
	})

	t.Run("Invalid Request body",func(t *testing.T) {
		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		app := fiber.New()
		app.Patch("/user/order/:id/status", testMiddleware,orderHandler.UpdateOrderStatusByUser)

		reqBody := []byte(`Invalid`)

		req := httptest.NewRequest("PATCH", "/user/order/1/status",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request body")
	})

	t.Run("Unauthorized",func(t *testing.T) {
		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", nil)
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
	

		app := fiber.New()
		app.Patch("/user/order/:id/status", testMiddleware,orderHandler.UpdateOrderStatusByUser)

		reqBody:= []byte(`{
			"status":"pending"
		}`)

		req := httptest.NewRequest("PATCH", "/user/order/1/status",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Unauthorized")
	})

	t.Run("UpdateOrderStatusByUser Success",func(t *testing.T) {
		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "failID")
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)

		app := fiber.New()
		app.Patch("/user/order/:id/status", testMiddleware,orderHandler.UpdateOrderStatusByUser)

		reqBody:= []byte(`{
			"status":"pending"
		}`)

		req := httptest.NewRequest("PATCH", "/user/order/1/status",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid user ID")
	})

	t.Run("UpdateOrderStatusByUser Success",func(t *testing.T) {
		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			Status: "pending",
		}

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		orderService.On("FindOrderById",mock.Anything).Return(&orderMock,nil)
		orderService.On("UpdateStatusByUser",mock.Anything,mock.Anything,mock.Anything).Return(errors.New("Update status failed"))

		app := fiber.New()
		app.Patch("/user/order/:id/status", testMiddleware,orderHandler.UpdateOrderStatusByUser)

		reqBody:= []byte(`{
			"status":"pending"
		}`)

		req := httptest.NewRequest("PATCH", "/user/order/1/status",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Update status failed")
	})

}

func TestGetAllOrderByUserId(t *testing.T) {
	t.Run("GetAllOrderByUserId Success",func(t *testing.T) {
		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			Status: "pending",
		}

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		orderService.On("GetAllOrderByUserId",mock.Anything).Return(&orderMock,nil)

		app := fiber.New()
		app.Get("/user/order", testMiddleware,orderHandler.GetAllOrderByUserId)

		req := httptest.NewRequest("GET", "/user/order",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Get All Orders Success")
	})

	t.Run("UserID is nil",func(t *testing.T) {

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", nil)
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		app := fiber.New()
		app.Get("/user/order", testMiddleware,orderHandler.GetAllOrderByUserId)

		req := httptest.NewRequest("GET", "/user/order",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Unauthorized")
	})

	t.Run("Invalid format UserID",func(t *testing.T) {

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "failUserID")
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		app := fiber.New()
		app.Get("/user/order", testMiddleware,orderHandler.GetAllOrderByUserId)

		req := httptest.NewRequest("GET", "/user/order",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid user ID format")
	})
}

func TestGetAllOrders(t *testing.T) {
	t.Run("GetAllOrders Success",func(t *testing.T) {
		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			Status: "pending",
		}
		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		orderService.On("GetAllOrdersAdmin").Return(&orderMock,nil)

		app := fiber.New()
		app.Get("/admin/order",orderHandler.GetAllOrders)

		req := httptest.NewRequest("GET", "/admin/order",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Get All Order Success")
	})
	t.Run("Error GetAllOrders",func(t *testing.T) {
		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		orderService.On("GetAllOrdersAdmin").Return(nil,errors.New("Get Order Error"))

		app := fiber.New()
		app.Get("/admin/order",orderHandler.GetAllOrders)

		req := httptest.NewRequest("GET", "/admin/order",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Get Order Error")
	})
}

func TestUpdateOrderStatusByAdmin(t *testing.T) {
	t.Run("UpdateOrderStatusByAdmin Success",func(t *testing.T) {
		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		orderService.On("UpdateStatusByAdmin",mock.Anything,mock.Anything).Return(nil)

		app := fiber.New()
		app.Patch("/admin/order/:id/status", testMiddleware,orderHandler.UpdateOrderStatusByAdmin)

		reqBody:= []byte(`{
			"status":"complete"
		}`)

		req := httptest.NewRequest("PATCH", "/admin/order/1/status",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Update status success")
	})
	t.Run("Invalid product ID",func(t *testing.T) {
		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		app := fiber.New()
		app.Patch("/admin/order/:id/status", testMiddleware,orderHandler.UpdateOrderStatusByAdmin)

		reqBody:= []byte(`{
			"status":"complete"
		}`)

		req := httptest.NewRequest("PATCH", "/admin/order/failID/status",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid product ID")
	})

	t.Run("Invalid request body",func(t *testing.T) {
		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		app := fiber.New()
		app.Patch("/admin/order/:id/status", testMiddleware,orderHandler.UpdateOrderStatusByAdmin)

		reqBody:= []byte(`failInvalid body`)

		req := httptest.NewRequest("PATCH", "/admin/order/1/status",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request body")
	})

	t.Run("Error UpdateOrderStatusByAdmin",func(t *testing.T) {
		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		orderService.On("UpdateStatusByAdmin",mock.Anything,mock.Anything).Return(errors.New("Update status failed"))

		app := fiber.New()
		app.Patch("/admin/order/:id/status", testMiddleware,orderHandler.UpdateOrderStatusByAdmin)

		reqBody:= []byte(`{
			"status":"complete"
		}`)

		req := httptest.NewRequest("PATCH", "/admin/order/1/status",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Update status failed")
	})
}

func TestDeleteOrder(t *testing.T) {
	t.Run("Delete Order Success",func(t *testing.T) {
		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		orderService.On("DeleteOrder",mock.Anything).Return(nil)

		app := fiber.New()
		app.Delete("/admin/order/:id", testMiddleware,orderHandler.DeleteOrder)

		req := httptest.NewRequest("DELETE", "/admin/order/1",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Order deleted successfully")
	})

	t.Run("Invalid order ID",func(t *testing.T) {
		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		orderService.On("DeleteOrder",mock.Anything).Return(nil)

		app := fiber.New()
		app.Delete("/admin/order/:id", testMiddleware,orderHandler.DeleteOrder)

		req := httptest.NewRequest("DELETE", "/admin/order/fail",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid order ID")
	})

	t.Run("Delete Order Success",func(t *testing.T) {
		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		orderService := services.NewOrderServiceMock()
		orderHandler := handlers.NewOrderHandler(orderService)
		
		orderService.On("DeleteOrder",mock.Anything).Return(errors.New("Error to delete"))

		app := fiber.New()
		app.Delete("/admin/order/:id", testMiddleware,orderHandler.DeleteOrder)

		req := httptest.NewRequest("DELETE", "/admin/order/1",nil)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Error to delete")
	})

}