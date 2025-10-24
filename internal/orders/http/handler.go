package http

import (
	"dunhayat-api/internal/auth/http"
	"dunhayat-api/internal/orders"
	"dunhayat-api/internal/orders/usecase"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	createOrderUseCase usecase.CreateOrderUseCase
}

func NewOrderHandler(
	createOrderUseCase usecase.CreateOrderUseCase,
) *OrderHandler {
	return &OrderHandler{
		createOrderUseCase: createOrderUseCase,
	}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	// Extract user ID from authentication context
	userID, ok := http.GetUserIDFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	var req orders.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	order, err := h.createOrderUseCase.Execute(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	response := fiber.Map{
		"message": "Order created successfully",
		"data":    order,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")
	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order ID is required",
		})
	}

	// TODO: Implement GetOrderUseCase

	response := fiber.Map{
		"message":  "Get order endpoint - to be implemented",
		"order_id": orderID,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *OrderHandler) CancelOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")
	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order ID is required",
		})
	}

	// TODO: Implement CancelOrderUseCase

	response := fiber.Map{
		"message":  "Cancel order endpoint - to be implemented",
		"order_id": orderID,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
