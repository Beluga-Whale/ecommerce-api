package services

import (
	"errors"

	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/repositories"
)

type ReviewServiceInterface interface {
	GetReviewsByUserID(userIDUint uint) ([]models.Review , error) 
	CreateReview(userIDUint uint, req dto.CreateReviewDTO) error
	GetReviewAll(productId uint ) (dto.ReviewAllProductSummaryResponse,error)
}

type ReviewService struct{
	reviewRepo repositories.ReviewRepositoryInterface
}

func NewReviewService(reviewRepo repositories.ReviewRepositoryInterface)*ReviewService {
	return &ReviewService{reviewRepo:reviewRepo}
}

func (s *ReviewService) GetReviewsByUserID(userIDUint uint) ([]models.Review , error) {
	reviews,err := s.reviewRepo.GetUserReviews(userIDUint)

	if err != nil {
		return nil,errors.New("Error to get user reviews")
	}

	return reviews,nil
}

func (s *ReviewService) CreateReview(userIDUint uint, req dto.CreateReviewDTO) error {
	hasPurchased, err := s.reviewRepo.HasPurchasedProduct(userIDUint, req.ProductID)
	if err != nil {
		return err
	}
	if !hasPurchased {
		return errors.New("you cannot review this product because you haven’t purchased it")
	}

	review := &models.Review{
		UserID:    userIDUint,
		ProductID: req.ProductID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}
	return s.reviewRepo.Create(review)
}

func (s *ReviewService) GetReviewAll(productId uint ) (dto.ReviewAllProductSummaryResponse,error) {
	var response dto.ReviewAllProductSummaryResponse

	// NOTE - Get all reviews
	reviews, err := s.reviewRepo.GetReviewAllByProductId(productId)
	if err != nil {
		return response, errors.New("Error to get all review Product")
	}

	// NOTE - Calculate average rating
	avg, err := s.reviewRepo.GetAverageRatingByProductId(productId)
	if err != nil {
		return response, errors.New("Error to get average rating")
	}

	// NOTE - Count per star
	countMap := make(map[int]int)
	for _, review := range reviews {
		countMap[int(review.Rating)]++
	}

	response = dto.ReviewAllProductSummaryResponse{
		Average:      avg,
		Total:        len(reviews),
		CountPerStar: countMap,
		ReviewList:   reviews,
	}

	return response, nil
}