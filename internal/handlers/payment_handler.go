package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/Beluga-Whale/ecommerce-api/internal/dto"
	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"
)

type PaymentHandlerInterface interface {
	CreatePaymentIntent(c *fiber.Ctx) error
}

type StripeHandler struct{
	orderService services.OrderServiceInterface
}

func NewStripeHandler(orderService services.OrderServiceInterface) *StripeHandler {
	return &StripeHandler{orderService:orderService}
}

func (h *StripeHandler) CreatePaymentIntent(c *fiber.Ctx) error {
	var req dto.CreatePaymentIntentRequestDTO

	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(req.Amount),
		Currency: stripe.String(string("usd")),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}
	params.AddMetadata("orderId", strconv.Itoa(int(req.OrderID)))
	params.AddMetadata("userId", strconv.Itoa(int(req.UserID)))

	intent, err := paymentintent.New(params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"clientSecret": intent.ClientSecret,
	})
}

func (h *StripeHandler) Webhook(c *fiber.Ctx) error {
	stripeKey := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if stripeKey == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "Missing Stripe webhook secret")
	}

	payload := c.Body()
	sigHeader := c.Get("Stripe-Signature")

	event, err := stripe.ConstructEvent(payload, sigHeader, stripeKey)
	if err != nil {
		fmt.Println(" Webhook signature verification failed:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Webhook Error")
	}

	switch event.Type {
case "payment_intent.succeeded":
    var paymentIntent stripe.PaymentIntent
    err := json.Unmarshal(event.Data.Raw, &paymentIntent)
    if err != nil {
        fmt.Println(" Failed to parse payment intent:", err)
        return c.Status(fiber.StatusBadRequest).SendString("Invalid payment data")
    }

    fmt.Println(" PaymentIntent was successful")

    orderID, _ := strconv.Atoi(paymentIntent.Metadata["orderId"])
    userID, _ := strconv.Atoi(paymentIntent.Metadata["userId"])
    status := models.Status("paid")
    orderIDUint := uint(orderID)

    err = h.orderService.UpdateStatusOrder(&orderIDUint, status, uint(userID))
    if err != nil {
        fmt.Println(" Failed to update order:", err)
        return c.Status(fiber.StatusInternalServerError).SendString("Update failed")
    }

    fmt.Println("✅ Order updated")
	case "payment_intent.payment_failed":
		fmt.Println(" Payment failed:", event.ID)
	default:
		fmt.Println("Unhandled event type:", event.Type)
	}

	return c.SendStatus(fiber.StatusOK)
}