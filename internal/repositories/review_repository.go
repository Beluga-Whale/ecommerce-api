package repositories

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type ReviewRepositoryInterface interface{
	GetUserReviews(userIDUint uint) ([]models.Review, error)
	Create(review *models.Review ) error 
	HasPurchasedProduct(userID uint, productID uint) (bool, error)
	GetReviewAllByProductId(productId uint) ([]dto.ReviewAllProduct,error)
}

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository{
	return &ReviewRepository{db:db}
}

func (r *ReviewRepository) GetUserReviews(userIDUint uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("user_id = ?", userIDUint).Find(&reviews).Error
	return reviews, err
}

func (r *ReviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewRepository)HasPurchasedProduct(userID uint, productID uint) (bool, error){
	var count int64

	err := r.db.
		Table("orders").
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Joins("JOIN product_variants ON order_items.product_variant_id = product_variants.id").
		Where("orders.user_id = ? AND product_variants.product_id = ? AND orders.status = ?", userID, productID, "complete").
		Count(&count).Error

	return count > 0, err
}

func (r *ReviewRepository)GetReviewAllByProductId(productId uint) ([]dto.ReviewAllProduct,error){
	var reviews []dto.ReviewAllProduct

	err := r.db.Table("products").
		Joins("JOIN reviews ON products.id = reviews.product_id").
		Joins("JOIN users ON  reviews.user_id = users.id").
		Where("products.id = ? ",productId).
		Scan(&reviews).Error

	if err != nil {
		return nil, err
	}
	return reviews,nil
	
}