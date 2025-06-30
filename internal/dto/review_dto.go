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

type ReviewAllProduct struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	ProductID uint   `json:"productId"`
	Rating    int64  `json:"rating"`
	Comment   string `json:"comment"`
	Avatar    string `json:"avatar"`
}

type ReviewAllProductSummaryResponse struct {
	Average      float64            `json:"average"`
	Total        int                `json:"total"`
	CountPerStar map[int]int        `json:"countPerStar"` // 1–5 stars
	ReviewList   []ReviewAllProduct `json:"reviewList"`
}