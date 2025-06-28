package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/repositories"
	"github.com/Beluga-Whale/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type OrderServiceInterface interface {
	CreateOrder(userID uint, req dto.CreateOrderRequestDTO) (*models.Order, error)
	CancelOrderAndRestoreStock( orderID uint) error
	UpdateStatusOrder(orderID *uint, status models.Status,userId uint) error
	GetOrderByID(orderID uint, userIDUint uint) (*models.Order, error)
	GetAllOrderByUserId(userIDUint uint) ([]models.Order,error)
	UpdateStatusByUser(userIDUint uint,orderID *uint, status models.Status) error
	GetAllOrdersAdmin() ([]models.Order,error)
	UpdateStatusByAdmin(orderID *uint, status models.Status) error
	GetDashboardSummary() (*dto.DashboardSummaryDTO, error)
	GetProductTop() ([]dto.TopProductDTO,error)
}

type OrderService struct {
	db 		    *gorm.DB
	orderRepo   repositories.OrderRepositoryInterface
	productUtil utils.ProductInterface
}

func NewOrderService(db *gorm.DB,orderRepo repositories.OrderRepositoryInterface,productUtil utils.ProductInterface) *OrderService {
	return &OrderService{
		db: db,
		orderRepo: orderRepo,
		productUtil:productUtil,
	}
}


func (s *OrderService) CreateOrder(userID uint, req dto.CreateOrderRequestDTO) (*models.Order, error) {
	// NOTE - เช็ค ว่า req มีมาจิรงไหม
    if len(req.Items) == 0 {
		return nil, errors.New("no item in order")
	}

	// NOTE - เก็บ productVariantID เป็น slice []
	productVariantIDs := []uint{}
	for _,item := range req.Items {
		productVariantIDs = append(productVariantIDs, item.VariantID)
	}

	productVariants,err := s.orderRepo.FindProductVariantByID(productVariantIDs)

	if err != nil {
		return nil,errors.New("fail to find product by productID")
	}

	// NOTE - Create Transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// NOTE - recover
	defer func() {
		if r:= recover(); r!= nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var total float64
	orderItems := []models.OrderItem{}

	for _,item := range req.Items{
		productV := s.productUtil.FindProductVariantID(productVariants, item.VariantID)
		if productV == nil {
			return nil, errors.New("productVariant not found")
		}

		if productV.Stock < int(item.Quantity){
			return nil,errors.New("stock not enough")
		}

		// // NOTE - ตัด stock
		productV.Stock -= int(item.Quantity)

		if err := s.orderRepo.UpdateProductVariantStock(tx,productV.ID,productV.Stock); err !=nil {
			tx.Rollback()
			return nil, errors.New("failed to update product stock")
		}

		// NOTE - คิดเงินรวม
		total += (productV.Price - *productV.Product.SalePrice) * float64(item.Quantity)

		orderItems = append(orderItems, models.OrderItem{
			ProductVariantID: productV.ID,
			Quantity: item.Quantity,
			PriceAtPurchase: productV.Price,
		})
	}

	order := models.Order{
		UserID: 	userID,
		FullName:   req.FullName,
		Phone: 		req.Phone,
		Address:  	req.Address,
		Province:   req.Province,
		District:   req.District,
		Subdistrict: req.Subdistrict,
		Zipcode:    req.Zipcode,
		TotalPrice: total,
		Status: 	models.Pending,
		OrderItem: 	orderItems,
		PaymentExpireAt: time.Now().Add(10 *time.Second),
	}

	// NOTE - create Order


	if err := s.orderRepo.Create(tx,&order); err != nil {
		tx.Rollback()
		return nil,err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	orderWithProducts, err := s.orderRepo.FindByIDWithItemsAndProducts(order.ID)
	if err != nil {
		return nil, err
	}

	return orderWithProducts,nil
}

func (s *OrderService) CancelOrderAndRestoreStock( orderID uint) error {
	var order models.Order

	if err := s.db.Preload("OrderItem").First(&order, orderID).Error; err != nil {
		return err
	}

	if order.Status == models.Pending {
		for _, item := range order.OrderItem {
			// NOTE - หา order จาก orderID จากนั้นจะทำการ update stock ตาม Quantity ของ oderID ตัวนั้นๆ
			if err := s.db.Model(&models.ProductVariant{}).
				Where("id = ?",item.ProductVariantID).
				Update("stock", gorm.Expr("stock + ?",item.Quantity)).Error; err != nil {
					return err
				}
		}
		// NOTE -เปลี่ยนสถานะเป็น cancel
		if err := s.db.Model(&order).Update("status",models.Cancel).Error; err != nil{
			return err
		}
	}
	return nil
}

func (s *OrderService) UpdateStatusOrder(orderID *uint, status models.Status,userId uint) error {

	if orderID == nil {
		return errors.New("no order id")
	}

	// NOTE - เช็คว่า orderID มีค่าไหม
	order, err := s.orderRepo.FindOrderById(*orderID)
	if err != nil {
		return  fmt.Errorf("orderRepo.FindByIDWithItemsAndProducts failed: %w", err)
	}

	if order == nil {
		return errors.New("order not found")
	}

	// NOTE - เช็คว่า userID ตรงกับ order.UserID ไหม
	if order.UserID != userId {
		return errors.New("unauthorized to update this order")
	}

	if err := s.orderRepo.UpdateStatusOrder(orderID,status); err != nil {
		return fmt.Errorf("orderRepo.UpdateStatusOrder failed: %w", err)
	}

	return nil
}

func (s *OrderService) GetOrderByID(orderID uint, userIDUint uint) (*models.Order, error) {
	order, err := s.orderRepo.FindOrderById(orderID)
	if err != nil {
		return nil, fmt.Errorf("orderRepo.FindByIDWithItemsAndProducts failed: %w", err)
	}

	if order == nil {
		return nil, errors.New("order not found")
	}

	// NOTE - เช็คว่า userID ตรงกับ order.UserID ไหม
	if order.UserID != userIDUint {
		return nil, errors.New("unauthorized to update this order")
	}


	return order, nil
}

func (s *OrderService) GetAllOrderByUserId(userIDUint uint) ([]models.Order,error) {

	orderAll, err := s.orderRepo.FindAllOrderByUserId(userIDUint)

	if err != nil {
		return nil,errors.New("Error to find to order")
	}

	return orderAll,nil
}
func (s *OrderService) UpdateStatusByUser(userIDUint uint,orderID *uint, status models.Status) error {
	if orderID == nil {
		return errors.New("no order id")
	}

	// NOTE - เช็คว่า orderID มีค่าไหม
	order, err := s.orderRepo.FindOrderById(*orderID)
	if err != nil {
		return  fmt.Errorf("orderRepo.FindByIDWithItemsAndProducts failed: %w", err)
	}

	if order == nil {
		return errors.New("order not found")
	}

	// NOTE - เช็คว่า userID ตรงกับ order.UserID ไหม
	if order.UserID != userIDUint {
		return  errors.New("unauthorized to update this order")
	}

	if err = s.orderRepo.UpdateStatusOrderByUserId(uint(*orderID),status); err !=nil{
		return errors.New("Order can not update status")
	}

	return nil
	
}

func (s *OrderService) GetAllOrdersAdmin() ([]models.Order,error) {
	orders,err := s.orderRepo.FindAll()

	if err != nil {
		return nil,err
	}

	return orders,nil
}

func (s *OrderService) UpdateStatusByAdmin(orderID *uint, status models.Status) error {
	if orderID == nil {
		return errors.New("no order id")
	}

	// NOTE - เช็คว่า orderID มีค่าไหม
	order, err := s.orderRepo.FindOrderById(*orderID)
	if err != nil {
		return  fmt.Errorf("orderRepo.FindByIDWithItemsAndProducts failed: %w", err)
	}

	if order == nil {
		return errors.New("order not found")
	}

	if err = s.orderRepo.UpdateStatusOrderByUserId(uint(*orderID),status); err !=nil{
		return errors.New("Order can not update status")
	}

	return nil
}

func percentDiff(current, previous float64) float64 {
	if previous == 0 {
		if current > 0 {
			return 100
		}
		return 0
	}
	return ((current - previous) / previous) * 100
}

func (s *OrderService) GetDashboardSummary() (*dto.DashboardSummaryDTO, error) {
	now := time.Now()
	startOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfLastMonth := startOfThisMonth.AddDate(0, -1, 0)
	endOfLastMonth := startOfThisMonth.Add(-time.Nanosecond)

	var orderTotal int64
	var ordersThisMonth int64
	var ordersLastMonth int64
	var revenueThisMonth float64
	var revenueLastMonth float64
	var customersThisMonth int64
	var customersLastMonth int64
	var statusPending int64
	var statusPaid int64
	var statusShipped int64
	var statusCancel int64


	// NOTE - Status
	s.db.Model(&models.Order{}).
		Where("status = ?","pending").
		Count(&statusPending)

	s.db.Model(&models.Order{}).
		Where("status = ?","paid").
		Count(&statusPaid)

	s.db.Model(&models.Order{}).
		Where("status = ?","shipped").
		Count(&statusShipped)

	s.db.Model(&models.Order{}).
		Where("status = ?","cancel").
		Count(&statusCancel)

	//NOTE - Orders

	s.db.Model(&models.Order{}).Count(&orderTotal)

	s.db.Model(&models.Order{}).
		Where("created_at >= ?", startOfThisMonth).
		Count(&ordersThisMonth)

	s.db.Model(&models.Order{}).
		Where("created_at >= ? AND created_at <= ?", startOfLastMonth, endOfLastMonth).
		Count(&ordersLastMonth)

	//NOTE - Revenue
	s.db.Model(&models.Order{}).
		Where("created_at >= ?", startOfThisMonth).
		Select("SUM(total_price)").Scan(&revenueThisMonth)

	s.db.Model(&models.Order{}).
		Where("created_at >= ? AND created_at <= ?", startOfLastMonth, endOfLastMonth).
		Select("SUM(total_price)").Scan(&revenueLastMonth)

	//NOTE - Customers
	s.db.Model(&models.Order{}).
		Where("created_at >= ?", startOfThisMonth).
		Distinct("user_id").Count(&customersThisMonth)

	s.db.Model(&models.Order{}).
		Where("created_at >= ? AND created_at <= ?", startOfLastMonth, endOfLastMonth).
		Distinct("user_id").Count(&customersLastMonth)


	//NOTE - Growth percent
	orderGrowth := percentDiff(float64(ordersThisMonth), float64(ordersLastMonth))
	revenueGrowth := percentDiff(revenueThisMonth, revenueLastMonth)
	customerGrowth := percentDiff(float64(customersThisMonth), float64(customersLastMonth))

	summary := &dto.DashboardSummaryDTO{
		OrderTotal: 		   int(orderTotal),
		OrdersThisMonth:       int(ordersThisMonth),
		OrdersLastMonth:       int(ordersLastMonth),
		OrderGrowthPercent:    orderGrowth,
		RevenueThisMonth:      revenueThisMonth,
		RevenueLastMonth:      revenueLastMonth,
		RevenueGrowthPercent:  revenueGrowth,
		CustomersThisMonth:    int(customersThisMonth),
		CustomersLastMonth:    int(customersLastMonth),
		CustomerGrowthPercent: customerGrowth,
		StatusPending: 		   int(statusPending),
		StatusPaid:            int(statusPaid),
		StatusShipped:         int(statusShipped),
		StatusCancel:          int(statusCancel),
	}
	
	return summary, nil
}

func (s *OrderService) GetProductTop() ([]dto.TopProductDTO,error) {
	topProduct,err := s.orderRepo.GetTop5ProductsBySales()
	if err != nil {
		return nil,errors.New("Error to query top product")
	}

	return topProduct,nil
}
