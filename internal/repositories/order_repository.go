package repositories

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	FindProductByID(productIDs []uint)([]models.Product,error)
	Create(order *models.Order) error
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository{
	return &OrderRepository{db:db}
}

func (r *OrderRepository) FindProductByID(productIDs []uint)([]models.Product,error) {
	var products []models.Product

	err := r.db.Where("id IN ?", productIDs).Find(&products).Error

	if err != nil {
		return nil,err
	}

	return products,nil 
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}


		if err := tx.Preload("OrderItem.Product").First(order, order.ID).Error; err != nil {
			return err
		}
		return nil
	})
}