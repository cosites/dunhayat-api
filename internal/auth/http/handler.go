package http

import (
	"dunhayat-api/internal/auth"
	"dunhayat-api/internal/auth/usecase"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	requestOTPUseCase usecase.RequestOTPUseCase
	verifyOTPUseCase  usecase.VerifyOTPUseCase
}

func NewAuthHandler(
	requestOTPUseCase usecase.RequestOTPUseCase,
	verifyOTPUseCase usecase.VerifyOTPUseCase,
) *AuthHandler {
	return &AuthHandler{
		requestOTPUseCase: requestOTPUseCase,
		verifyOTPUseCase:  verifyOTPUseCase,
	}
}

func (h *AuthHandler) RequestOTP(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPost {
		return c.Status(
			fiber.StatusMethodNotAllowed,
		).JSON(fiber.Map{
			"error": "Method not allowed",
		})
	}

	var req auth.RequestOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Phone number is required",
		})
	}

	otp, err := h.requestOTPUseCase.Execute(c.Context(), req.Phone)
	if err != nil {
		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(fiber.Map{
			"error": "Failed to send OTP",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "OTP sent successfully",
		"phone":      otp.Phone,
		"expires_at": otp.ExpiresAt,
	})
}

func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPost {
		return c.Status(
			fiber.StatusMethodNotAllowed,
		).JSON(fiber.Map{
			"error": "Method not allowed",
		})
	}

	var req auth.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Phone == "" || req.Code == "" {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"error": "Phone number and code are required",
		})
	}

	authResponse, err := h.verifyOTPUseCase.Execute(
		c.Context(),
		req.Phone,
		req.Code,
	)
	if err != nil {
		switch err.Error() {
		case "OTP has expired":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "OTP has expired",
			})
		case "invalid OTP code":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid OTP code",
			})
		default:
			return c.Status(
				fiber.StatusInternalServerError,
			).JSON(fiber.Map{
				"error": "Verification failed",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": "OTP verified successfully",
			"user":    authResponse.User,
			"token":   authResponse.Token,
		},
	)
}

func (h *AuthHandler) GetOTPStatus(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return c.Status(
			fiber.StatusMethodNotAllowed,
		).JSON(fiber.Map{
			"error": "Method not allowed",
		})
	}

	phone := c.Query("phone")
	if phone == "" {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(fiber.Map{
			"error": "Phone number is required",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"phone":  phone,
		"status": "pending",
	})
}
