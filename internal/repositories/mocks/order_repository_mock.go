package repositories

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type OrderRepositoryMock struct {
	mock.Mock
}

func NewOrderRepositoryMock() *OrderRepositoryMock {
	return &OrderRepositoryMock{}
}


func (m *OrderRepositoryMock)FindProductVariantByID(productVariantIDs []uint)([]models.ProductVariant,error){
	args := m.Called(productVariantIDs)

	if orders,ok := args.Get(0).([]models.ProductVariant);ok{
		return orders,nil
	}
	return nil, args.Error(1)
}

func (m *OrderRepositoryMock)Create(tx *gorm.DB,order *models.Order) error {
	args := m.Called(tx,order)
	return args.Error(0)
}

func (m *OrderRepositoryMock)UpdateProductVariantStock(tx *gorm.DB,productVariantID uint, newStock int) error {
	args := m.Called(tx,productVariantID,newStock)
	return args.Error(0)
}

func (m *OrderRepositoryMock)FindByIDWithItemsAndProducts(orderID uint) (*models.Order, error) {
	args := m.Called(orderID)
	if order,ok := args.Get(0).(*models.Order);ok {
		return order,nil
	}
	return nil,args.Error(1)
}

func (m *OrderRepositoryMock)UpdateStatusOrder(orderId *uint, status models.Status) error {
	args := m.Called(orderId,status)
	return args.Error(0)
}

func (m *OrderRepositoryMock)FindOrderById(orderID uint) (*models.Order, error){
	args := m.Called(orderID)
	if order,ok := args.Get(0).(*models.Order);ok {
		return order,nil
	}
	return nil,args.Error(1)
}

func (m *OrderRepositoryMock)FindAllOrderByUserId(userIDUint uint) ([]models.Order,error){
	args := m.Called(userIDUint)
	if order,ok := args.Get(0).([]models.Order);ok {
		return order,nil
	}
	return nil,args.Error(1)
}


func (m *OrderRepositoryMock)UpdateStatusOrderByUserId(orderID uint,status models.Status) error{
	args := m.Called(orderID,status)
	return args.Error(0)
}

func (m *OrderRepositoryMock)FindAll() ([]models.Order,error) {
	args := m.Called()
	if order,ok := args.Get(0).([]models.Order);ok {
		return order,nil
	}
	return nil,args.Error(1)
}

func (m *OrderRepositoryMock)GetTop5ProductsBySales() ([]dto.TopProductDTO, error){
	args := m.Called()
	if topProduct,ok := args.Get(0).([]dto.TopProductDTO);ok {
		return topProduct,nil
	}
	return nil,args.Error(1)
}
func (m *OrderRepositoryMock)GetSalesPerDay() ([]dto.SalesPerMonthDTO, error) {
	args := m.Called()
	if salePerMonth,ok := args.Get(0).([]dto.SalesPerMonthDTO);ok {
		return salePerMonth,nil
	}
	return nil,args.Error(1)
}

func (m *OrderRepositoryMock)Delete(id uint) error{
	args := m.Called(id)
	return args.Error(0)
}

func (m *OrderRepositoryMock)GetUserDetail() ([]dto.CustomerDTO,error) {
	args := m.Called()
	if customer,ok := args.Get(0).([]dto.CustomerDTO);ok {
		return customer,nil
	}
	return nil,args.Error(1)
}