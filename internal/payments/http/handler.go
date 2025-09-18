package http

import (
	"dunhayat-api/internal/payments"
	"dunhayat-api/internal/payments/usecase"

	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	initiatePaymentUseCase  usecase.InitiatePaymentUseCase
	verifyPaymentUseCase    usecase.VerifyPaymentUseCase
	handleCallbackUseCase   usecase.HandleCallbackUseCase
	getPaymentStatusUseCase usecase.GetPaymentStatusUseCase
}

func NewPaymentHandler(
	initiatePaymentUseCase usecase.InitiatePaymentUseCase,
	verifyPaymentUseCase usecase.VerifyPaymentUseCase,
	handleCallbackUseCase usecase.HandleCallbackUseCase,
	getPaymentStatusUseCase usecase.GetPaymentStatusUseCase,
) *PaymentHandler {
	return &PaymentHandler{
		initiatePaymentUseCase:  initiatePaymentUseCase,
		verifyPaymentUseCase:    verifyPaymentUseCase,
		handleCallbackUseCase:   handleCallbackUseCase,
		getPaymentStatusUseCase: getPaymentStatusUseCase,
	}
}

func (h *PaymentHandler) InitiatePayment(c *fiber.Ctx) error {
	var req payments.InitiatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	response, err := h.initiatePaymentUseCase.Execute(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Payment initiated successfully",
		"data":    response,
	})
}

func (h *PaymentHandler) VerifyPayment(c *fiber.Ctx) error {
	var req payments.VerifyPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	response, err := h.verifyPaymentUseCase.Execute(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment verification completed",
		"data":    response,
	})
}

func (h *PaymentHandler) HandleCallback(c *fiber.Ctx) error {
	var callbackData payments.PaymentCallbackRequest
	if err := c.BodyParser(&callbackData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid callback data: " + err.Error(),
		})
	}

	err := h.handleCallbackUseCase.Execute(c.Context(), callbackData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Callback processed successfully",
	})
}

func (h *PaymentHandler) GetPaymentStatus(c *fiber.Ctx) error {
	orderID := c.Query("order_id")
	trackingCode := c.Query("tracking_code")

	if orderID == "" && trackingCode == "" {
		orderID = c.Params("id")
	}

	if orderID == "" && trackingCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Either order_id or tracking_code must be provided",
		})
	}

	req := &payments.GetPaymentStatusRequest{
		OrderID:      orderID,
		TrackingCode: trackingCode,
	}

	response, err := h.getPaymentStatusUseCase.Execute(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment status retrieved successfully",
		"data":    response,
	})
}
