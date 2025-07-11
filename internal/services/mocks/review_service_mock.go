package services

import (
	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type ReviewServiceMock struct {
	mock.Mock
}

func NewReviewServiceMock() *ReviewServiceMock {
	return &ReviewServiceMock{}
}

func (m *ReviewServiceMock) GetReviewsByUserID(userIDUint uint) ([]models.Review , error) {
	args := m.Called(userIDUint)
	if reviews,ok := args.Get(0).([]models.Review);ok {
		return reviews,args.Error(1)
	}
	return nil,args.Error(1)
}

func (m *ReviewServiceMock)CreateReview(userIDUint uint, req dto.CreateReviewDTO) error {
	args := m.Called(userIDUint,req)

	return args.Error(0)
}

func (m *ReviewServiceMock) GetReviewAll(productId uint ) (dto.ReviewAllProductSummaryResponse,error) {
	args := m.Called(productId)
	if reviews,ok := args.Get(0).(dto.ReviewAllProductSummaryResponse);ok {
		return reviews,args.Error(1)
	}
	return dto.ReviewAllProductSummaryResponse{},args.Error(1)
}