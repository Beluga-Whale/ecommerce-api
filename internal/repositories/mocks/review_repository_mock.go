package repositories

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type ReviewRepositoryMock struct {
	mock.Mock
}

func NewReviewRepositoryMock() *ReviewRepositoryMock{
	return &ReviewRepositoryMock{}
}

func (m *ReviewRepositoryMock) GetUserReviews(userIDUint uint) ([]models.Review, error) {
	args := m.Called(userIDUint)
	if userReview,ok := args.Get(0).([]models.Review);ok {
		return userReview,args.Error(1)
	}
	return nil,args.Error(1)
}


func (m *ReviewRepositoryMock) Create(review *models.Review ) error  {
	args := m.Called(review)
	
	return args.Error(0)
}

func (m *ReviewRepositoryMock) HasPurchasedProduct(userID uint, productID uint) (bool, error) {
	args := m.Called(userID,productID)
	return args.Bool(0),args.Error(1)
}


func (m *ReviewRepositoryMock) GetReviewAllByProductId(productId uint) ([]dto.ReviewAllProduct,error) {
	args := m.Called(productId)
	if reviewByProductID,ok := args.Get(0).([]dto.ReviewAllProduct);ok {
		return reviewByProductID,args.Error(1)
	}
	return nil,args.Error(1)
}


func (m *ReviewRepositoryMock) GetAverageRatingByProductId(productId uint) (float64, error) {
	args := m.Called(productId)

	return args.Get(0).(float64),args.Error(1)
}

