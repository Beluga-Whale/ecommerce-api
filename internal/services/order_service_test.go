package services_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	repositories "github.com/Beluga-Whale/ecommerce-api/internal/repositories/mocks"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	utils "github.com/Beluga-Whale/ecommerce-api/internal/utils/mocks"
	sqliteDriver "github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitializeDB(t *testing.T) *gorm.DB {
  	db, err := gorm.Open(sqliteDriver.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to open database: %v", err))
	}
	err = db.AutoMigrate(&models.Order{}, &models.OrderItem{}, &models.ProductVariant{}, &models.Product{})
	if err != nil {
		t.Fatalf("failed to auto migrate: %v", err)
	}
  return db
}

func TestCreateOrder(t *testing.T) {

	t.Run("Create order success",func(t *testing.T) {
		
		salePrice := 10.0
		productMock := &models.Product{SalePrice: &salePrice}
		
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		variantMock := models.ProductVariant{
			Model:gorm.Model{
				ID: 1,
			},
			Stock: 10,
			Size: "S",
			Price: 100.0,
			Product: *productMock,
		}

		mockOrder := &models.Order{
		Model: gorm.Model{ID: 1},
		UserID: 1,
		OrderItem: []models.OrderItem{
			{
				ProductVariantID: 1,
				Quantity:         2,
				PriceAtPurchase:  100,
			},
		},
		TotalPrice: 180,
		}

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		orderRepo.On("FindProductVariantByID",mock.Anything).Return([]models.ProductVariant{variantMock})
		productUtil.On("FindProductVariantID", mock.Anything, mock.Anything).Return(&variantMock)

		orderRepo.On("UpdateProductVariantStock",mock.Anything,mock.Anything,mock.Anything).Return(nil)
		orderRepo.On("Create",mock.Anything,mock.Anything).Return(nil)
		orderRepo.On("FindByIDWithItemsAndProducts",mock.Anything).Return(mockOrder,nil)

		req :=dto.CreateOrderRequestDTO{
			FullName:    "John Doe",
			Phone:       "0988888888",
			Address:     "123 Main St",
			Province:    "Bangkok",
			District:    "District",
			Subdistrict: "Subdistrict",
			Zipcode:     "10200",
			Items: []dto.CreateOrderItemDTO{
				{
					VariantID: 1,
					Quantity:  2,
				},
			},
		}

		order,err := orderService.CreateOrder(1,req)

		assert.NoError(t,err)
		assert.NotNil(t,order)
		assert.Equal(t,uint(1),order.UserID)
		assert.Equal(t,180.0,order.TotalPrice)

		productUtil.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
	})

	t.Run("Req items is zero",func(t *testing.T) {
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		req :=dto.CreateOrderRequestDTO{
			FullName:    "John Doe",
			Phone:       "0988888888",
			Address:     "123 Main St",
			Province:    "Bangkok",
			District:    "District",
			Subdistrict: "Subdistrict",
			Zipcode:     "10200",
			Items: []dto.CreateOrderItemDTO{},
		}

		_,err := orderService.CreateOrder(1,req)

		assert.EqualError(t,err,"no item in order")
	})

	t.Run("Error to find product variantById",func(t *testing.T) {
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		orderRepo.On("FindProductVariantByID",mock.Anything).Return(nil,errors.New("fail to find product by productID"))

		req :=dto.CreateOrderRequestDTO{
			FullName:    "John Doe",
			Phone:       "0988888888",
			Address:     "123 Main St",
			Province:    "Bangkok",
			District:    "District",
			Subdistrict: "Subdistrict",
			Zipcode:     "10200",
			Items: []dto.CreateOrderItemDTO{
				{
					VariantID: 1,
					Quantity:  2,
				},
			},
		}

		order,err := orderService.CreateOrder(1,req)

		assert.EqualError(t,err,"fail to find product by productID")
		assert.Nil(t,order)
		productUtil.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
	})

	t.Run("Error to update productVariant",func(t *testing.T) {
		
		salePrice := 10.0
		productMock := &models.Product{SalePrice: &salePrice}
		
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		variantMock := models.ProductVariant{
			Model:gorm.Model{
				ID: 1,
			},
			Stock: 10,
			Size: "S",
			Price: 100.0,
			Product: *productMock,
		}
		

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		orderRepo.On("FindProductVariantByID",mock.Anything).Return([]models.ProductVariant{variantMock})
		productUtil.On("FindProductVariantID", mock.Anything, mock.Anything).Return(&variantMock)

		orderRepo.On("UpdateProductVariantStock",mock.Anything,mock.Anything,mock.Anything).Return(errors.New("failed to update product stock"))
		

		req :=dto.CreateOrderRequestDTO{
			FullName:    "John Doe",
			Phone:       "0988888888",
			Address:     "123 Main St",
			Province:    "Bangkok",
			District:    "District",
			Subdistrict: "Subdistrict",
			Zipcode:     "10200",
			Items: []dto.CreateOrderItemDTO{
				{
					VariantID: 1,
					Quantity:  2,
				},
			},
		}

		order,err := orderService.CreateOrder(1,req)

		assert.EqualError(t,err,"failed to update product stock")
		assert.Nil(t,order)

		productUtil.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
	})
	
	t.Run("Error to create order",func(t *testing.T) {
		
		salePrice := 10.0
		productMock := &models.Product{SalePrice: &salePrice}
		
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		variantMock := models.ProductVariant{
			Model:gorm.Model{
				ID: 1,
			},
			Stock: 10,
			Size: "S",
			Price: 100.0,
			Product: *productMock,
		}

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		orderRepo.On("FindProductVariantByID",mock.Anything).Return([]models.ProductVariant{variantMock})
		productUtil.On("FindProductVariantID", mock.Anything, mock.Anything).Return(&variantMock)

		orderRepo.On("UpdateProductVariantStock",mock.Anything,mock.Anything,mock.Anything).Return(nil)
		orderRepo.On("Create",mock.Anything,mock.Anything).Return(errors.New("Error to create order"))


		req :=dto.CreateOrderRequestDTO{
			FullName:    "John Doe",
			Phone:       "0988888888",
			Address:     "123 Main St",
			Province:    "Bangkok",
			District:    "District",
			Subdistrict: "Subdistrict",
			Zipcode:     "10200",
			Items: []dto.CreateOrderItemDTO{
				{
					VariantID: 1,
					Quantity:  2,
				},
			},
		}

		order,err := orderService.CreateOrder(1,req)

		assert.EqualError(t,err,"Error to create order")
		assert.Nil(t,order)

		productUtil.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
	})

	t.Run("Error to findOrderById",func(t *testing.T) {
		
		salePrice := 10.0
		productMock := &models.Product{SalePrice: &salePrice}
		
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		variantMock := models.ProductVariant{
			Model:gorm.Model{
				ID: 1,
			},
			Stock: 10,
			Size: "S",
			Price: 100.0,
			Product: *productMock,
		}

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		orderRepo.On("FindProductVariantByID",mock.Anything).Return([]models.ProductVariant{variantMock})
		productUtil.On("FindProductVariantID", mock.Anything, mock.Anything).Return(&variantMock)

		orderRepo.On("UpdateProductVariantStock",mock.Anything,mock.Anything,mock.Anything).Return(nil)
		orderRepo.On("Create",mock.Anything,mock.Anything).Return(nil)
		orderRepo.On("FindByIDWithItemsAndProducts",mock.Anything).Return(nil,errors.New("Error to find order by id"))

		req :=dto.CreateOrderRequestDTO{
			FullName:    "John Doe",
			Phone:       "0988888888",
			Address:     "123 Main St",
			Province:    "Bangkok",
			District:    "District",
			Subdistrict: "Subdistrict",
			Zipcode:     "10200",
			Items: []dto.CreateOrderItemDTO{
				{
					VariantID: 1,
					Quantity:  2,
				},
			},
		}

		order,err := orderService.CreateOrder(1,req)

		assert.EqualError(t,err,"Error to find order by id")
		assert.Nil(t,order)

		productUtil.AssertExpectations(t)
		orderRepo.AssertExpectations(t)
	})
}

func TestUpdateStatusOrder(t *testing.T) {
	t.Run("UpdateStatusOrder Success",func(t *testing.T) {
		orderID := uint(1)

		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			Status: "pending",
		}

		orderRepo.On("FindOrderById",orderID).Return(&orderMock,nil)
		orderRepo.On("UpdateStatusOrder",&orderID,models.Status("pending")).Return(nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.UpdateStatusOrder(&orderID,"pending",1)

		assert.NoError(t,err)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Not have orderID",func(t *testing.T) {
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.UpdateStatusOrder(nil,"pending",1)

		assert.Error(t,err)
		assert.EqualError(t,err,"no order id")

		orderRepo.AssertExpectations(t)
	})

	t.Run("Error to FindOrderById",func(t *testing.T) {
		orderID := uint(1)

		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()
		orderRepo.On("FindOrderById",mock.Anything).Return(nil,errors.New("orderRepo.FindByIDWithItemsAndProducts failed"))

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.UpdateStatusOrder(&orderID,"pending",1)

		assert.Error(t,err)
		assert.EqualError(t,err,"orderRepo.FindByIDWithItemsAndProducts failed: orderRepo.FindByIDWithItemsAndProducts failed")
		orderRepo.AssertExpectations(t)
	})

	t.Run("Order is Nil",func(t *testing.T) {
		orderID := uint(1)

		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()
		orderRepo.On("FindOrderById",orderID).Return(nil,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.UpdateStatusOrder(&orderID,"pending",1)

		assert.Error(t,err)
		assert.EqualError(t,err,"order not found")
		orderRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized to update this order",func(t *testing.T) {
		orderID := uint(1)

		mockOrder := models.Order{
			UserID: 2,
		}

		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()
		orderRepo.On("FindOrderById",orderID).Return(&mockOrder,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.UpdateStatusOrder(&orderID,"pending",1)

		assert.Error(t,err)
		assert.EqualError(t,err,"unauthorized to update this order")
		orderRepo.AssertExpectations(t)
	})

	t.Run("Error to update Status order",func(t *testing.T) {
		orderID := uint(1)

		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			Status: "pending",
		}

		orderRepo.On("FindOrderById",orderID).Return(&orderMock,nil)
		orderRepo.On("UpdateStatusOrder",&orderID,models.Status("pending")).Return(errors.New("orderRepo.UpdateStatusOrder failed"))

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.UpdateStatusOrder(&orderID,"pending",1)

		assert.Error(t,err)
		assert.EqualError(t,err,"orderRepo.UpdateStatusOrder failed: orderRepo.UpdateStatusOrder failed")

		orderRepo.AssertExpectations(t)
	})
}

func TestGetOrderByID(t *testing.T) {
	t.Run("GetOrderByID Success",func(t *testing.T) {
		orderId := uint(1)
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()
		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			Status: "pending",
			TotalPrice: 100.0,
			UserID: 1,
		}


		orderRepo.On("FindOrderById",orderId).Return(&orderMock,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		orders,err := orderService.GetOrderByID(1,1)

		assert.NoError(t,err)

		assert.Contains(t,"pending",orders.Status)

		orderRepo.AssertExpectations(t)

	})

	t.Run("Error to findOrderById",func(t *testing.T) {
		orderId := uint(1)
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindOrderById",orderId).Return(nil,errors.New("orderRepo.FindByIDWithItemsAndProducts failed:"))

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		_,err := orderService.GetOrderByID(1,1)

		assert.EqualError(t,err,"orderRepo.FindByIDWithItemsAndProducts failed: orderRepo.FindByIDWithItemsAndProducts failed:")

		orderRepo.AssertExpectations(t)

	})

	t.Run("Order Not found",func(t *testing.T) {
		orderId := uint(1)
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindOrderById",orderId).Return(nil,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		_,err := orderService.GetOrderByID(1,1)

		assert.EqualError(t,err,"order not found")


		orderRepo.AssertExpectations(t)

	})

	t.Run("Unauthorized to update order",func(t *testing.T) {
		orderId := uint(1)
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()
		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			Status: "pending",
			TotalPrice: 100.0,
			UserID: 2,
		}


		orderRepo.On("FindOrderById",orderId).Return(&orderMock,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		_,err := orderService.GetOrderByID(1,1)

		assert.EqualError(t,err,"unauthorized to update this order")


		orderRepo.AssertExpectations(t)
	})
}

func TestGetAllOrderByUserId(t *testing.T) {
	t.Run("GetAllOrderByUserId Success",func(t *testing.T) {
		userId := uint(1)
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderAll := []models.Order{
			{
				Model:      gorm.Model{ID: 1},
				Status:     "pending",
				TotalPrice: 100.0,
				UserID:     1,
			},
			{
				Model:      gorm.Model{ID: 2},
				Status:     "paid",
				TotalPrice: 200.0,
				UserID:     1,
			},
		}

		orderRepo.On("FindAllOrderByUserId",userId).Return(orderAll,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		orders,err :=orderService.GetAllOrderByUserId(uint(1))

		assert.NoError(t,err)
		assert.Len(t,orders,2)
		assert.Equal(t, models.Status("pending") , orders[0].Status)
		assert.Equal(t, models.Status("paid"), orders[1].Status)

		orderRepo.AssertExpectations(t)
	})
	t.Run("Error to findAllOrder",func(t *testing.T) {
		userId := uint(1)
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindAllOrderByUserId",userId).Return(nil,errors.New("Error to find to order"))

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		_,err :=orderService.GetAllOrderByUserId(uint(1))

		assert.EqualError(t,err,"Error to find to order")

		orderRepo.AssertExpectations(t)
	})
}

func TestUpdateStatusByUser(t *testing.T) {
	t.Run("UpdateStatus ByUser ID Success",func(t *testing.T) {
		orderId := uint(1)
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			Status: "pending",
		}

		orderRepo.On("FindOrderById",orderId).Return(&orderMock,nil)
		orderRepo.On("UpdateStatusOrderByUserId",orderId,mock.Anything).Return(nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.UpdateStatusByUser(1,&orderId,models.Status("pending"))
		assert.NoError(t,err)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Order Id Is Nil",func(t *testing.T) {
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.UpdateStatusByUser(1,nil,models.Status("pending"))

		assert.EqualError(t,err,"no order id")
		orderRepo.AssertExpectations(t)
	})

	t.Run("Error to findOrderByID",func(t *testing.T) {
		orderId := uint(1)
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindOrderById",orderId).Return(nil,errors.New("orderRepo.FindByIDWithItemsAndProducts failed"))

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		
		err := orderService.UpdateStatusByUser(1,&orderId,models.Status("pending"))

		assert.EqualError(t,err,"orderRepo.FindByIDWithItemsAndProducts failed: orderRepo.FindByIDWithItemsAndProducts failed")
		orderRepo.AssertExpectations(t)
	})

	t.Run("UpdateStatus ByUser ID Success",func(t *testing.T) {
		orderId := uint(1)
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindOrderById",orderId).Return(nil,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.UpdateStatusByUser(1,&orderId,models.Status("pending"))

		assert.EqualError(t,err,"order not found")

		orderRepo.AssertExpectations(t)
	})

	t.Run("UserId not equal UserID in order",func(t *testing.T) {
		orderId := uint(1)
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			Status: "pending",
		}

		orderRepo.On("FindOrderById",orderId).Return(&orderMock,nil)
		
		orderService := services.NewOrderService(db,orderRepo,productUtil)
		
		err := orderService.UpdateStatusByUser(2,&orderId,models.Status("pending"))
		assert.EqualError(t,err,"unauthorized to update this order")

		orderRepo.AssertExpectations(t)
	})

	t.Run("Error to update status order",func(t *testing.T) {
		orderId := uint(1)
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			Status: "pending",
		}

		orderRepo.On("FindOrderById",orderId).Return(&orderMock,nil)
		orderRepo.On("UpdateStatusOrderByUserId",orderId,models.Status("pending")).Return(errors.New("Order can not update status"))

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.UpdateStatusByUser(1,&orderId,models.Status("pending"))
		assert.EqualError(t,err,"Order can not update status")

		orderRepo.AssertExpectations(t)
	})
}

func TestGetAllOrdersAdmin(t *testing.T) {
	t.Run("GetAllOderAdmin Success",func(t *testing.T) {
		orderMock := []models.Order{{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			TotalPrice: 100.0,},
		}

		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindAll").Return(orderMock,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		orders,err := orderService.GetAllOrdersAdmin()

		assert.NoError(t,err)
		assert.Equal(t,uint(1),orders[0].UserID)
	})
	t.Run("Error to getAllOrderAdmin",func(t *testing.T) {
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindAll").Return(nil,errors.New("Error to find order"))

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		_,err := orderService.GetAllOrdersAdmin()

		assert.EqualError(t,err,"Error to find order")
	})
}

func TestUpdateStatusByAdmin(t *testing.T) {
	t.Run("UpdateStatusByAdmin Success",func(t *testing.T) {
		orderId := uint(1)
		status := models.Status("paid")

		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			Status: "complete",
		}

		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindOrderById",orderId).Return(&orderMock,nil)
		orderRepo.On("UpdateStatusOrderByUserId",orderId,status).Return(nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.UpdateStatusByAdmin(&orderId,status)

		assert.NoError(t,err)

		orderRepo.AssertExpectations(t)
	})
	t.Run("OrderID is nil",func(t *testing.T) {
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		status := models.Status("paid")
		err := orderService.UpdateStatusByAdmin(nil,status)

		assert.EqualError(t,err,"no order id")

		orderRepo.AssertExpectations(t)
	})

	t.Run("Error to findOrderByID",func(t *testing.T) {
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindOrderById",mock.Anything).Return(nil,errors.New("orderRepo.FindByIDWithItemsAndProducts failed"))

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		orderId := uint(1)
		status := models.Status("paid")
		err := orderService.UpdateStatusByAdmin(&orderId,status)

		assert.EqualError(t,err,"orderRepo.FindByIDWithItemsAndProducts failed: orderRepo.FindByIDWithItemsAndProducts failed")

		orderRepo.AssertExpectations(t)
	})

	t.Run("Order is nil",func(t *testing.T) {
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindOrderById",mock.Anything).Return(nil,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		orderId := uint(1)
		status := models.Status("paid")
		err := orderService.UpdateStatusByAdmin(&orderId,status)

		assert.EqualError(t,err,"order not found")

		orderRepo.AssertExpectations(t)
	})

	t.Run("Error to update status",func(t *testing.T) {
		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			Status: "complete",
		}

		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindOrderById",mock.Anything).Return(&orderMock,nil)
		orderRepo.On("UpdateStatusOrderByUserId",mock.Anything,mock.Anything).Return(errors.New("Order can not update status"))

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		orderId := uint(1)
		status := models.Status("paid")
		err := orderService.UpdateStatusByAdmin(&orderId,status)

		assert.EqualError(t,err,"Order can not update status")
		orderRepo.AssertExpectations(t)
	})
}

func TestGetProductTop(t *testing.T) {
	t.Run("GetDashboardSummary Success",func(t *testing.T) {
		topProduct := []dto.TopProductDTO{
			{
				ProductID: 1,
				Name: "A",
				TotalSold: 100.0,
			},
			{
				ProductID: 2,
				Name: "B",
				TotalSold: 200.0,
			},
			{
				ProductID: 3,
				Name: "C",
				TotalSold: 100.0,
			},
			{
				ProductID: 4,
				Name: "D",
				TotalSold: 200.0,
			},
			{
				ProductID: 5,
				Name: "E",
				TotalSold: 200.0,
			},
		}
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("GetTop5ProductsBySales").Return(topProduct,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		top5Product,err := orderService.GetProductTop()

		assert.NoError(t,err)

		assert.Equal(t,top5Product[0].Name,"A")

		orderRepo.AssertExpectations(t)
	})

	t.Run("GetDashboardSummary Success",func(t *testing.T) {
	
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("GetTop5ProductsBySales").Return(nil,errors.New("Error to query top product"))

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		top5Product,err := orderService.GetProductTop()

		assert.EqualError(t,err,"Error to query top product")

		assert.Nil(t,top5Product)

		orderRepo.AssertExpectations(t)
	})
}

func TestGetSalesChartData(t *testing.T) {
	t.Run("GetSalesChartData Success",func(t *testing.T) {
		salesPerDay := []dto.SalesPerMonthDTO{
			{
				Date:      time.Date(2025, 7, 11, 0, 0, 0, 0, time.UTC),
				TotalSale: 1000.0,
			},
			{
				Date:      time.Date(2025, 7, 12, 0, 0, 0, 0, time.UTC),
				TotalSale: 2000.0,
			},
		}

		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("GetSalesPerDay").Return(salesPerDay,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		result,err := orderService.GetSalesChartData()

		assert.NoError(t,err)

		assert.Equal(t,result[0].TotalSale,1000.0)
		assert.Equal(t,result[1].TotalSale,2000.0)

		orderRepo.AssertExpectations(t)
	})
	t.Run("Error GetSalesChartData",func(t *testing.T) {

		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("GetSalesPerDay").Return(nil,errors.New("Error to query salePreDay"))

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		result,err := orderService.GetSalesChartData()

		assert.EqualError(t,err,"Error to query salePreDay")
		assert.Nil(t,result)

		orderRepo.AssertExpectations(t)
	})
}

func TestDeleteOrder(t *testing.T) {
	t.Run("DeleteOrder Success",func(t *testing.T) {
		id := uint(1)
		existingOrder := models.Order{
			Model: gorm.Model{ID: 1},
		}

		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindOrderById",id).Return(&existingOrder,nil)
		orderRepo.On("Delete",id).Return(nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.DeleteOrder(id)

		assert.NoError(t,err)

		orderRepo.AssertExpectations(t)
	})

	t.Run("Error FindOrderById",func(t *testing.T) {
		id := uint(1)
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindOrderById",id).Return(nil,errors.New("Error finding product"))

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.DeleteOrder(id)

		assert.EqualError(t,err,"Error finding product")

		orderRepo.AssertExpectations(t)
	})

	t.Run("Order is Nil",func(t *testing.T) {
		id := uint(1)
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindOrderById",id).Return(nil,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.DeleteOrder(id)

		assert.EqualError(t,err,"Product not found")

		orderRepo.AssertExpectations(t)
	})

	t.Run("Error to Delete Order",func(t *testing.T) {
		id := uint(1)
		existingOrder := models.Order{
			Model: gorm.Model{ID: 1},
		}

		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindOrderById",id).Return(&existingOrder,nil)
		orderRepo.On("Delete",id).Return(errors.New("Error deleting order"))

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.DeleteOrder(id)

		assert.EqualError(t,err,"Error deleting order")

		orderRepo.AssertExpectations(t)
	})
}

func TestGetCustomerDetail(t *testing.T) {
	t.Run("GetCustomerDetail Success",func(t *testing.T) {
		customers := []dto.CustomerDTO{
			{
				ID: 1,
				Name: "TEST A",
			},
			{
				ID: 2,
				Name: "TEST B",
			},
		}

		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("GetUserDetail").Return(customers,nil)
		
		orderService := services.NewOrderService(db,orderRepo,productUtil)

		result,err := orderService.GetCustomerDetail()

		assert.NoError(t,err)
		assert.Equal(t,result[0].Name,"TEST A")
		assert.Equal(t,result[1].Name,"TEST B")
	})

	t.Run("Error To GetCustomerDetail",func(t *testing.T) {
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("GetUserDetail").Return(nil,errors.New("Error to query customer detail"))
		
		orderService := services.NewOrderService(db,orderRepo,productUtil)

		result,err := orderService.GetCustomerDetail()

		assert.EqualError(t,err,"Error to query customer detail")
		assert.Nil(t,result)
	})
}