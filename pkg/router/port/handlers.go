package port

import (
	"github.com/gofiber/fiber/v2"
)

type ProductHandler interface {
	ListProducts(c *fiber.Ctx) error
	GetProduct(c *fiber.Ctx) error
}

type OrderHandler interface {
	CreateOrder(c *fiber.Ctx) error
	GetOrder(c *fiber.Ctx) error
	CancelOrder(c *fiber.Ctx) error
}

type AuthHandler interface {
	RequestOTP(c *fiber.Ctx) error
	VerifyOTP(c *fiber.Ctx) error
	GetOTPStatus(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type AuthMiddleware interface {
	Authenticate() fiber.Handler
}
