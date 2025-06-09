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
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository{
	return &OrderRepository{db:db}
}

func (r *OrderRepository) FindProductVariantByID(productVariantIDs []uint)([]models.ProductVariant,error) {
	var productVariants []models.ProductVariant

	err := r.db.Where("id IN ?", productVariantIDs).Find(&productVariants).Error

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

	if	err := r.db.Preload("OrderItem.ProductVariant.Product").First(&order,orderID).Error; err != nil {
		return nil,err
	}

	return &order,nil

}