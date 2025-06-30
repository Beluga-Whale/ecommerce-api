package dto

type CreateReviewDTO struct {
	ProductID uint   `json:"productId",required`
	Rating    int64  `json:"rating",required`
	Comment   string `json:"comment",required`
}

type ReviewResponse struct {
	ProductID uint   `json:"productId"`
	Rating    int64  `json:"rating"`
	Comment   string `json:"comment"`
}
