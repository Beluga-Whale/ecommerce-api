package services

import (
	"errors"
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
		total += productV.Price * float64(item.Quantity)

		orderItems = append(orderItems, models.OrderItem{
			ProductVariantID: productV.ID,
			Quantity: item.Quantity,
			PriceAtPurchase: productV.Price,
		})
	}

	order := models.Order{
		UserID: 	userID,
		Phone: 		req.Phone,
		Address:  	req.Address,
		Note: 		req.Note,
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