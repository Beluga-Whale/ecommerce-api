package repositories

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	FindProductVariantByID(productVariantIDs []uint)([]models.ProductVariant,error)
	Create(tx *gorm.DB,order *models.Order) error
	UpdateProductVariantStock(tx *gorm.DB,productVariantID uint, newStock int) error
	FindByIDWithItemsAndProducts(orderID uint) (*models.Order, error)
	UpdateStatusOrder(orderId *uint, status models.Status) error
	FindOrderById(orderID uint) (*models.Order, error)
	FindAllOrderByUserId(userIDUint uint) ([]models.Order,error)
	UpdateStatusOrderByUserId(orderID uint,status models.Status) error
	FindAll() ([]models.Order,error)
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository{
	return &OrderRepository{db:db}
}

func (r *OrderRepository) FindProductVariantByID(productVariantIDs []uint)([]models.ProductVariant,error) {
	var productVariants []models.ProductVariant

	err := r.db.Preload("Product").Where("id IN ?", productVariantIDs).Find(&productVariants).Error

	if err != nil {
		return nil,err
	}

	return productVariants,nil 
}

func (r *OrderRepository) Create(tx *gorm.DB,order *models.Order) error {
	return tx.Create(order).Error
}

func (r *OrderRepository) UpdateProductVariantStock(tx *gorm.DB,productVariantID uint, newStock int) error {
	return tx.Model(&models.ProductVariant{}).Where("id = ?",productVariantID).Update("stock",newStock).Error
}

func (r *OrderRepository) FindByIDWithItemsAndProducts(orderID uint) (*models.Order, error) {
	var order models.Order

	if	err := r.db.Preload("OrderItem.ProductVariant.Product").Preload("User").Preload("Coupon").First(&order,orderID).Error; err != nil {
		return nil,err
	}

	return &order,nil

}

func (r *OrderRepository) UpdateStatusOrder(orderId *uint, status models.Status) error {
	if err := r.db.Model(&models.Order{}).Where("id = ?",*orderId).Update("status",status).Error; err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) FindOrderById(orderID uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("User").Preload("Coupon").Preload("OrderItem.ProductVariant.Product").Where("id = ?", orderID).First(&order).Error

	if err != nil {
		return nil, err
	}
	
	return &order, nil
}

func (r *OrderRepository) FindAllOrderByUserId(userIDUint uint) ([]models.Order,error) {
	var orderAll []models.Order

	err := r.db.Preload("Coupon").Preload("OrderItem.ProductVariant.Product").Where("user_id = ?",userIDUint).Order("id DESC").Find(&orderAll).Error
	if err != nil {
		return nil,err
	}

	return orderAll,nil
}

func (r *OrderRepository) UpdateStatusOrderByUserId(orderID uint,status models.Status) error{
	err := r.db.Model(&models.Order{}).Where("id = ?",orderID).Update("status",status).Error

	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) FindAll() ([]models.Order,error){
	var orders []models.Order
	
	err := r.db.Preload("Coupon").Preload("OrderItem.ProductVariant.Product").Find(&orders).Error

	return orders,err
}