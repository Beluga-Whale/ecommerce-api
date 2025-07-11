package handlers_test

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/handlers"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	servicesMock "github.com/Beluga-Whale/ecommerce-api/internal/services/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetUserReviews(t *testing.T) {
	t.Run("GetUserReviews Success",func(t *testing.T) {
		reviews := []models.Review{
			{
				Model: gorm.Model{ID: 1},
				Rating: 4,
			},
			{
				Model: gorm.Model{ID: 2},
				Rating: 3,
			},
		}

		reviewService := servicesMock.NewReviewServiceMock()

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		reviewService.On("GetReviewsByUserID",uint(1)).Return(reviews,nil)

		reviewHandler :=  handlers.NewReviewHandler(reviewService)

		app := fiber.New()
		app.Get("/user/review",testMiddleware,reviewHandler.GetUserReviews)

		req :=httptest.NewRequest("GET","/user/review",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "User reviews")

	})

	t.Run("Unauthorized UserID Nil",func(t *testing.T) {
		reviews := []models.Review{
			{
				Model: gorm.Model{ID: 1},
				Rating: 4,
			},
			{
				Model: gorm.Model{ID: 2},
				Rating: 3,
			},
		}

		reviewService := servicesMock.NewReviewServiceMock()

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", nil)
			return c.Next()
		}

		reviewService.On("GetReviewsByUserID",uint(1)).Return(reviews,nil)

		reviewHandler :=  handlers.NewReviewHandler(reviewService)

		app := fiber.New()
		app.Get("/user/review",testMiddleware,reviewHandler.GetUserReviews)

		req :=httptest.NewRequest("GET","/user/review",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Unauthorized")

	})

	t.Run("Invalid user ID format",func(t *testing.T) {
		reviews := []models.Review{
			{
				Model: gorm.Model{ID: 1},
				Rating: 4,
			},
			{
				Model: gorm.Model{ID: 2},
				Rating: 3,
			},
		}

		reviewService := servicesMock.NewReviewServiceMock()

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "failToUserID")
			return c.Next()
		}

		reviewService.On("GetReviewsByUserID",uint(1)).Return(reviews,nil)

		reviewHandler :=  handlers.NewReviewHandler(reviewService)

		app := fiber.New()
		app.Get("/user/review",testMiddleware,reviewHandler.GetUserReviews)

		req :=httptest.NewRequest("GET","/user/review",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid user ID format")

	})

	t.Run("Error to get reviews",func(t *testing.T) {
		reviewService := servicesMock.NewReviewServiceMock()

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		reviewService.On("GetReviewsByUserID",uint(1)).Return(nil,errors.New("Failed to get reviews"))

		reviewHandler :=  handlers.NewReviewHandler(reviewService)

		app := fiber.New()
		app.Get("/user/review",testMiddleware,reviewHandler.GetUserReviews)

		req :=httptest.NewRequest("GET","/user/review",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Failed to get reviews")

	})
}

func TestCreateReviews(t *testing.T) {
	t.Run("CreateReview Success",func(t *testing.T) {
		reqMock := dto.CreateReviewDTO{
			ProductID: 1,
			Rating: 4,
			Comment: "GOOD",
		}

		reviewService := servicesMock.NewReviewServiceMock()

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		reviewService.On("CreateReview",uint(1),reqMock).Return(nil)

		reviewHandler :=  handlers.NewReviewHandler(reviewService)

		app := fiber.New()
		app.Post("/user/review",testMiddleware,reviewHandler.CreateReviews)

		reqBody:= []byte(`{
			"productId":1,
			"rating":4,
			"comment":"GOOD"
		}`)

		req :=httptest.NewRequest("POST","/user/review",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Review created successfully")

	})

	t.Run("Invalid Request Body",func(t *testing.T) {
		reviewService := servicesMock.NewReviewServiceMock()

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		reviewHandler :=  handlers.NewReviewHandler(reviewService)

		app := fiber.New()
		app.Post("/user/review",testMiddleware,reviewHandler.CreateReviews)

		reqBody:= []byte(`Invalid Request Body`)

		req :=httptest.NewRequest("POST","/user/review",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid request body")

	})

	t.Run("Unauthorized",func(t *testing.T) {
		reviewService := servicesMock.NewReviewServiceMock()

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", nil)
			return c.Next()
		}


		reviewHandler :=  handlers.NewReviewHandler(reviewService)

		app := fiber.New()
		app.Post("/user/review",testMiddleware,reviewHandler.CreateReviews)

		reqBody:= []byte(`{
			"productId":1,
			"rating":4,
			"comment":"GOOD"
		}`)

		req :=httptest.NewRequest("POST","/user/review",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Unauthorized")

	})

	t.Run("Invalid user ID format",func(t *testing.T) {
		reviewService := servicesMock.NewReviewServiceMock()

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "fail Invalid")
			return c.Next()
		}

		reviewHandler :=  handlers.NewReviewHandler(reviewService)

		app := fiber.New()
		app.Post("/user/review",testMiddleware,reviewHandler.CreateReviews)

		reqBody:= []byte(`{
			"productId":1,
			"rating":4,
			"comment":"GOOD"
		}`)

		req :=httptest.NewRequest("POST","/user/review",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid user ID format")

	})

	t.Run("Error CreateReview",func(t *testing.T) {
		reqMock := dto.CreateReviewDTO{
			ProductID: 1,
			Rating: 4,
			Comment: "GOOD",
		}

		reviewService := servicesMock.NewReviewServiceMock()

		testMiddleware := func(c *fiber.Ctx) error {
			c.Locals("userID", "1")
			return c.Next()
		}

		reviewService.On("CreateReview",uint(1),reqMock).Return(errors.New("Error to createReview"))

		reviewHandler :=  handlers.NewReviewHandler(reviewService)

		app := fiber.New()
		app.Post("/user/review",testMiddleware,reviewHandler.CreateReviews)

		reqBody:= []byte(`{
			"productId":1,
			"rating":4,
			"comment":"GOOD"
		}`)

		req :=httptest.NewRequest("POST","/user/review",bytes.NewReader(reqBody))
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Error to createReview")

	})
}

func TestGetReviewProductAllByProductId(t *testing.T) {
	t.Run("GetReviewProductAllByProductId Success",func(t *testing.T) {
		reviewMock := dto.ReviewAllProductSummaryResponse{
			Average: 4,
			Total: 4,
			ReviewList: []dto.ReviewAllProduct{
				{
					FirstName: "A",
					Rating: 4,
				},
			},
		}
		reviewService := servicesMock.NewReviewServiceMock()

		reviewService.On("GetReviewAll",uint(1)).Return(reviewMock,nil)

		reviewHandler :=  handlers.NewReviewHandler(reviewService)

		app := fiber.New()
		app.Get("/product/review-all/:id",reviewHandler.GetReviewProductAllByProductId)

		req :=httptest.NewRequest("GET","/product/review-all/1",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Get Review Product successfully")
	})

	t.Run("Invalid Param Url",func(t *testing.T) {
		reviewService := servicesMock.NewReviewServiceMock()


		reviewHandler :=  handlers.NewReviewHandler(reviewService)

		app := fiber.New()
		app.Get("/product/review-all/:id",reviewHandler.GetReviewProductAllByProductId)

		req :=httptest.NewRequest("GET","/product/review-all/InvalidParam",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Invalid product ID")
	})

	t.Run("Error to get ReviewByProductID",func(t *testing.T) {
		reviewService := servicesMock.NewReviewServiceMock()

		reviewService.On("GetReviewAll",uint(1)).Return(nil,errors.New("Error to get reviewByProductID"))

		reviewHandler :=  handlers.NewReviewHandler(reviewService)

		app := fiber.New()
		app.Get("/product/review-all/:id",reviewHandler.GetReviewProductAllByProductId)

		req :=httptest.NewRequest("GET","/product/review-all/1",nil)
		req.Header.Set("Content-Type","application/json")

		res,err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)

		body, _ := io.ReadAll(res.Body)
		assert.Contains(t, string(body), "Error to get reviewByProductID")
	})
}