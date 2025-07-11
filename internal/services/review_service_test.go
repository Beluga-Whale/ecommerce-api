package services_test

import (
	"errors"
	"testing"

	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	repositories "github.com/Beluga-Whale/ecommerce-api/internal/repositories/mocks"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGetReviewsByUserID(t *testing.T){
	t.Run("GetReviewsByUserID Success",func(t *testing.T) {
		userID := uint(1)
		reviewMock := []models.Review{
			{
				Model: gorm.Model{ID: 1},
				Rating: 3.0,
				UserID: 1,
			},
			{
				Model: gorm.Model{ID: 2},
				Rating: 4.0,
				UserID: 2,
			},
		}

		reviewRepo := repositories.NewReviewRepositoryMock()

		reviewRepo.On("GetUserReviews",userID).Return(reviewMock,nil)

		reviewService:= services.NewReviewService(reviewRepo)

		review,err := reviewService.GetReviewsByUserID(userID)

		assert.NoError(t,err)
		assert.Equal(t, int64(3), review[0].Rating)
		assert.Equal(t, int64(4), review[1].Rating)

		reviewRepo.AssertExpectations(t)

	})

	t.Run("Error to GetReviewsByUserID ",func(t *testing.T) {
		userID := uint(1)

		reviewRepo := repositories.NewReviewRepositoryMock()

		reviewRepo.On("GetUserReviews",userID).Return(nil,errors.New("Error to get user reviews"))

		reviewService:= services.NewReviewService(reviewRepo)

		review,err := reviewService.GetReviewsByUserID(userID)

		assert.EqualError(t,err,"Error to get user reviews")
		assert.Nil(t,review)
		reviewRepo.AssertExpectations(t)
	})
}

func TestCreateReview(t *testing.T) {
	t.Run("CreateReview Success",func(t *testing.T) {
		userIDUint := uint(1)
		reqMock:=dto.CreateReviewDTO{
			ProductID: 1,
			Rating: 4,
			Comment: "GOOD",
		}

		reviewRepo := repositories.NewReviewRepositoryMock()

		reviewRepo.On("HasPurchasedProduct",userIDUint,reqMock.ProductID).Return(true,nil)
		reviewRepo.On("Create",mock.AnythingOfType("*models.Review")).Return(nil)

		reviewService:= services.NewReviewService(reviewRepo)

		err := reviewService.CreateReview(userIDUint,reqMock)

		assert.NoError(t,err)

		reviewRepo.AssertExpectations(t)

	})

	t.Run("Error to check",func(t *testing.T) {
		userIDUint := uint(1)
		reqMock:=dto.CreateReviewDTO{
			ProductID: 1,
			Rating: 4,
			Comment: "GOOD",
		}

		reviewRepo := repositories.NewReviewRepositoryMock()

		reviewRepo.On("HasPurchasedProduct",userIDUint,reqMock.ProductID).Return(false,errors.New("Error to check"))

		reviewService:= services.NewReviewService(reviewRepo)

		err := reviewService.CreateReview(userIDUint,reqMock)

		assert.EqualError(t,err,"Error to check")

		reviewRepo.AssertExpectations(t)

	})
	t.Run("It return false",func(t *testing.T) {
		userIDUint := uint(1)
		reqMock:=dto.CreateReviewDTO{
			ProductID: 1,
			Rating: 4,
			Comment: "GOOD",
		}

		reviewRepo := repositories.NewReviewRepositoryMock()

		reviewRepo.On("HasPurchasedProduct",userIDUint,reqMock.ProductID).Return(false,nil)

		reviewService:= services.NewReviewService(reviewRepo)

		err := reviewService.CreateReview(userIDUint,reqMock)

		assert.EqualError(t,err,"you cannot review this product because you haven’t purchased it")

		reviewRepo.AssertExpectations(t)

	})
}

func TestGetReviewAll(t *testing.T) {
	t.Run("GetReviewAll Success",func(t *testing.T) {
		productID := uint(1)
		reviewMock := []dto.ReviewAllProduct{
			{
				FirstName: "FirstName A",
				LastName: "LastName A",
				ProductID: 1,
				Rating: 4,
				Comment: "Comment A",

			},
		}
		reviewRepo := repositories.NewReviewRepositoryMock()

		reviewRepo.On("GetReviewAllByProductId",productID).Return(reviewMock,nil)
		reviewRepo.On("GetAverageRatingByProductId",productID).Return(4.0,nil)

		reviewService:= services.NewReviewService(reviewRepo)

		_,err := reviewService.GetReviewAll(productID)

		assert.NoError(t,err)

		reviewRepo.AssertExpectations(t)
	})

	t.Run("Error GetReviewAllByProductId",func(t *testing.T) {
		productID := uint(1)

		reviewRepo := repositories.NewReviewRepositoryMock()

		reviewRepo.On("GetReviewAllByProductId",productID).Return(nil,errors.New("Error to get all review Product"))

		reviewService:= services.NewReviewService(reviewRepo)

		_,err := reviewService.GetReviewAll(productID)

		assert.EqualError(t,err,"Error to get all review Product")

		reviewRepo.AssertExpectations(t)
	})

	t.Run("GetReviewAll Success",func(t *testing.T) {
		productID := uint(1)
		reviewMock := []dto.ReviewAllProduct{
			{
				FirstName: "FirstName A",
				LastName: "LastName A",
				ProductID: 1,
				Rating: 4,
				Comment: "Comment A",

			},
		}
		reviewRepo := repositories.NewReviewRepositoryMock()

		reviewRepo.On("GetReviewAllByProductId",productID).Return(reviewMock,nil)
		reviewRepo.On("GetAverageRatingByProductId",productID).Return(0.0,errors.New("Error to get average rating"))

		reviewService:= services.NewReviewService(reviewRepo)

		_,err := reviewService.GetReviewAll(productID)

		assert.EqualError(t,err,"Error to get average rating")

		reviewRepo.AssertExpectations(t)
	})
}