package repositories

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	FindProductByID(productIDs []uint)([]models.Product,error)
	Create(tx *gorm.DB,order *models.Order) error
	UpdateProductStock(tx *gorm.DB,productID uint, newStock int) error
	FindByIDWithItemsAndProducts(orderID uint) (*models.Order, error)
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

func (r *OrderRepository) Create(tx *gorm.DB,order *models.Order) error {
	return tx.Create(order).Error
}

func (r *OrderRepository) UpdateProductStock(tx *gorm.DB,productID uint, newStock int) error {
	return tx.Model(&models.Product{}).Where("id = ?",productID).Update("stock",newStock).Error
}

func (r *OrderRepository) FindByIDWithItemsAndProducts(orderID uint) (*models.Order, error) {
	var order models.Order

	if	err := r.db.Preload("OrderItem.Product").First(&order,orderID).Error; err != nil {
		return nil,err
	}

	return &order,nil

}