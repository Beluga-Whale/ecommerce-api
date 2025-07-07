package services

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type OrderServiceMock struct {
	mock.Mock
}

func NewOrderServiceMock() *OrderServiceMock {
	return &OrderServiceMock{}
}


func (m *OrderServiceMock) CreateOrder(userID uint, req dto.CreateOrderRequestDTO) (*models.Order, error){
	args := m.Called(userID,req)

	if order,ok := args.Get(0).(*models.Order); ok {
		return order,args.Error(1)
	}
	return nil,args.Error(1)
}

func (m *OrderServiceMock) CancelOrderAndRestoreStock( orderID uint) error{
	args := m.Called(orderID)

	return args.Error(1)
}


func (m *OrderServiceMock) UpdateStatusOrder(orderID *uint, status models.Status,userId uint) error{
	args := m.Called(orderID,status,userId)

	return args.Error(1)
}

func (m *OrderServiceMock) GetOrderByID(orderID uint, userIDUint uint) (*models.Order, error) {
	args := m.Called(orderID,userIDUint)

	if order,ok := args.Get(0).(*models.Order); ok {
		return order,args.Error(1)
	}
	return nil,args.Error(1)
}

func (m *OrderServiceMock) GetAllOrderByUserId(userIDUint uint) ([]models.Order,error) {
	args := m.Called(userIDUint)

	if order,ok := args.Get(0).([]models.Order); ok {
		return order,args.Error(1)
	}
	return nil,args.Error(1)
}

func (m *OrderServiceMock) UpdateStatusByUser(userIDUint uint,orderID *uint, status models.Status) error {
	args := m.Called(userIDUint,orderID,status)

	return args.Error(1)
}

func (m *OrderServiceMock) GetAllOrdersAdmin() ([]models.Order,error) {
	args := m.Called()

	if order,ok := args.Get(0).([]models.Order) ; ok {
		return order,args.Error(1)
	}

	return nil,args.Error(1)
}


func (m *OrderServiceMock) UpdateStatusByAdmin(orderID *uint, status models.Status) error {
	args := m.Called(orderID,status)

	return args.Error(1)
}

func (m *OrderServiceMock) GetDashboardSummary() (*dto.DashboardSummaryDTO, error) {
	args := m.Called()

	if summary,ok := args.Get(0).(dto.DashboardSummaryDTO);ok{
		return &summary,args.Error(1)
	}

	return nil,args.Error(1)
}

func (m *OrderServiceMock) GetProductTop() ([]dto.TopProductDTO,error) {
	args := m.Called()

	if topProduct,ok := args.Get(0).([]dto.TopProductDTO);ok{
		return topProduct,args.Error(1)
	}
	return nil,args.Error(1)
}

func (m *OrderServiceMock) GetSalesChartData() ([]dto.SalesPerMonthDTO, error) {
	args := m.Called()

	if salePerMonth,ok := args.Get(0).([]dto.SalesPerMonthDTO);ok{
		return salePerMonth,args.Error(1)
	}
	return nil,args.Error(1)
}

func (m *OrderServiceMock)DeleteOrder(id uint) error {
	args := m.Called(id)

	return args.Error(1)
}

func (m *OrderServiceMock)GetCustomerDetail() ([]dto.CustomerDTO,error) {
	args := m.Called()

	if customer,ok := args.Get(0).([]dto.CustomerDTO);ok{
		return customer,args.Error(1)
	}
	return nil,args.Error(1)
}

