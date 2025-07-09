package services_test

import (
	"errors"
	"fmt"
	"testing"

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

		orderRepo.On("FindOrderById",mock.Anything).Return(&orderMock,nil)
		orderRepo.On("UpdateStatusOrder",mock.Anything,mock.Anything).Return(nil)

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
		orderRepo.On("FindOrderById",mock.Anything).Return(nil,nil)

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
		orderRepo.On("FindOrderById",mock.Anything).Return(&mockOrder,nil)

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

		orderRepo.On("FindOrderById",mock.Anything).Return(&orderMock,nil)
		orderRepo.On("UpdateStatusOrder",mock.Anything,mock.Anything).Return(errors.New("orderRepo.UpdateStatusOrder failed"))

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		err := orderService.UpdateStatusOrder(&orderID,"pending",1)

		assert.Error(t,err)
		assert.EqualError(t,err,"orderRepo.UpdateStatusOrder failed: orderRepo.UpdateStatusOrder failed")

		orderRepo.AssertExpectations(t)
	})
}

func TestGetOrderByID(t *testing.T) {
	t.Run("GetOrderByID Success",func(t *testing.T) {
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()
		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			Status: "pending",
			TotalPrice: 100.0,
			UserID: 1,
		}


		orderRepo.On("FindOrderById",mock.Anything).Return(&orderMock,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		orders,err := orderService.GetOrderByID(1,1)

		assert.NoError(t,err)

		assert.Contains(t,"pending",orders.Status)

		orderRepo.AssertExpectations(t)

	})

	t.Run("Error to findOrderById",func(t *testing.T) {
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindOrderById",mock.Anything).Return(nil,errors.New("orderRepo.FindByIDWithItemsAndProducts failed:"))

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		_,err := orderService.GetOrderByID(1,1)

		assert.EqualError(t,err,"orderRepo.FindByIDWithItemsAndProducts failed: orderRepo.FindByIDWithItemsAndProducts failed:")

		orderRepo.AssertExpectations(t)

	})

	t.Run("Order Not found",func(t *testing.T) {
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindOrderById",mock.Anything).Return(nil,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		_,err := orderService.GetOrderByID(1,1)

		assert.EqualError(t,err,"order not found")


		orderRepo.AssertExpectations(t)

	})

	t.Run("Unauthorized to update order",func(t *testing.T) {
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()
		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			Status: "pending",
			TotalPrice: 100.0,
			UserID: 2,
		}


		orderRepo.On("FindOrderById",mock.Anything).Return(&orderMock,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		_,err := orderService.GetOrderByID(1,1)

		assert.EqualError(t,err,"unauthorized to update this order")


		orderRepo.AssertExpectations(t)
	})
}

func TestGetAllOrderByUserId(t *testing.T) {
	t.Run("GetAllOrderByUserId Success",func(t *testing.T) {
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

		orderRepo.On("FindAllOrderByUserId",mock.Anything).Return(orderAll,nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		orders,err :=orderService.GetAllOrderByUserId(uint(1))

		assert.NoError(t,err)
		assert.Len(t,orders,2)
		assert.Equal(t, models.Status("pending") , orders[0].Status)
		assert.Equal(t, models.Status("paid"), orders[1].Status)

		orderRepo.AssertExpectations(t)
	})
	t.Run("Error to findAllOrder",func(t *testing.T) {
		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindAllOrderByUserId",mock.Anything).Return(nil,errors.New("Error to find to order"))

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

		orderRepo.On("FindOrderById",mock.Anything).Return(&orderMock,nil)
		orderRepo.On("UpdateStatusOrderByUserId",mock.Anything,mock.Anything).Return(nil)

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

		orderRepo.On("FindOrderById",mock.Anything).Return(nil,errors.New("orderRepo.FindByIDWithItemsAndProducts failed"))

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

		orderRepo.On("FindOrderById",mock.Anything).Return(nil,nil)

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

		orderRepo.On("FindOrderById",mock.Anything).Return(&orderMock,nil)
		
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

		orderRepo.On("FindOrderById",mock.Anything).Return(&orderMock,nil)
		orderRepo.On("UpdateStatusOrderByUserId",mock.Anything,mock.Anything).Return(errors.New("Order can not update status"))

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
		orderMock := models.Order{
			Model: gorm.Model{ID: 1},
			UserID: 1,
			Status: "complete",
		}

		db  := InitializeDB(t)
		productUtil := utils.NewProductUtilMock()
		orderRepo := repositories.NewOrderRepositoryMock()

		orderRepo.On("FindOrderById",mock.Anything).Return(&orderMock,nil)
		orderRepo.On("UpdateStatusOrderByUserId",mock.Anything,mock.Anything).Return(nil)

		orderService := services.NewOrderService(db,orderRepo,productUtil)

		orderId := uint(1)
		status := models.Status("paid")
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