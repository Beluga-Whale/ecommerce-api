package services

import (
	"errors"

	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/repositories"
	"github.com/Beluga-Whale/ecommerce-api/internal/utils"
)

type OrderServiceInterface interface {
	CreateOrder(userID uint, req dto.CreateOrderRequestDTO) (*models.Order, error)
}

type OrderService struct {
	orderRepo repositories.OrderRepositoryInterface
	productUtil utils.ProductInterface
}

func NewOrderService(orderRepo repositories.OrderRepositoryInterface,productUtil utils.ProductInterface) *OrderService {
	return &OrderService{orderRepo: orderRepo,productUtil:productUtil }
}


func (s *OrderService) CreateOrder(userID uint, req dto.CreateOrderRequestDTO) (*models.Order, error) {
	// NOTE - เช็ค ว่า req มีมาจิรงไหม
    if len(req.Items) == 0 {
		return nil, errors.New("no item in order")
	}

	// NOTE - เก็บ productID เป็น slice []
	productIDs := []uint{}
	for _,item := range req.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	products,err := s.orderRepo.FindProductByID(productIDs)

	if err != nil {
		return nil,errors.New("fail to find product by productID")
	}

	// NOTE - Create Transaction
	tx := s.orderRepo.GetDB().Begin()
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
		product := s.productUtil.FindProductID(products, item.ProductID)
		if product == nil {
			return nil, errors.New("product not found")
		}

		if product.Stock < int(item.Quantity) {
			return nil,errors.New("stock not enough")
		}

		// NOTE - ตัด stock
		product.Stock -= int(item.Quantity)

		if err := s.orderRepo.UpdateProductStock(tx,product.ID,product.Stock); err !=nil {
			tx.Rollback()
			return nil, errors.New("failed to update product stock")
		}

		// NOTE - คิดเงินรวม
		total += product.Price * float64(item.Quantity)

		orderItems = append(orderItems, models.OrderItem{
			ProductID: product.ID,
			Quantity: item.Quantity,
			PriceAtPurchase: product.Price,
		})
	}

	order := models.Order{
		UserID: userID,
		Phone: req.Phone,
		Address:  req.Address,
		Note: req.Note,
		TotalPrice: total,
		Status: models.Pending,
		OrderItem: orderItems,
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