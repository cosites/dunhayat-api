package http

import (
	"strings"
	"time"

	"dunhayat-api/internal/auth"
	"dunhayat-api/internal/auth/port"
	"dunhayat-api/internal/auth/repository"
	"dunhayat-api/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	logger      *logger.Logger
	sessionRepo repository.SessionRepository
	userService port.UserReader
}

func NewAuthMiddleware(
	logger *logger.Logger,
	sessionRepo repository.SessionRepository,
	userService port.UserReader,
) *AuthMiddleware {
	return &AuthMiddleware{
		logger:      logger,
		sessionRepo: sessionRepo,
		userService: userService,
	}
}

func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		m.logger.Debug(
			"Authenticating request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("userAgent", c.Get("User-Agent")),
		)

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			m.logger.Warn(
				"Request rejected: missing authorization header",
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
			)
			return c.Status(
				fiber.StatusUnauthorized,
			).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			m.logger.Warn(
				"Request rejected: invalid authorization format",
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
			)
			return c.Status(
				fiber.StatusUnauthorized,
			).JSON(fiber.Map{
				"error": "Invalid authorization format. Use 'Bearer <token>'",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			m.logger.Warn(
				"Request rejected: empty token",
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
			)
			return c.Status(
				fiber.StatusUnauthorized,
			).JSON(fiber.Map{"error": "Token is required"})
		}

		session, err := m.sessionRepo.GetByToken(c.Context(), token)
		if err != nil {
			m.logger.Error(
				"Failed to validate session token",
				zap.Error(err),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
			)
			return c.Status(
				fiber.StatusUnauthorized,
			).JSON(fiber.Map{"error": "Invalid token"})
		}

		if session == nil {
			m.logger.Warn(
				"Request rejected: session not found",
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
			)
			return c.Status(
				fiber.StatusUnauthorized,
			).JSON(fiber.Map{"error": "Session not found"})
		}

		if time.Now().After(session.ExpiresAt) {
			m.logger.Warn(
				"Request rejected: session expired",
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.Time("expiresAt", session.ExpiresAt),
			)
			return c.Status(
				fiber.StatusUnauthorized,
			).JSON(fiber.Map{"error": "Session expired"})
		}

		user, err := m.userService.GetUserByID(
			c.Context(),
			session.UserID,
		)
		if err != nil {
			m.logger.Error(
				"Failed to get user information",
				zap.Error(err),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.String("userID", session.UserID.String()),
			)
			return c.Status(
				fiber.StatusInternalServerError,
			).JSON(fiber.Map{
				"error": "Failed to get user information",
			})
		}

		if user == nil {
			m.logger.Warn(
				"Request rejected: user not found",
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.String("userID", session.UserID.String()),
			)
			return c.Status(
				fiber.StatusUnauthorized,
			).JSON(fiber.Map{"error": "User not found"})
		}

		m.logger.Debug(
			"Request authenticated successfully",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("userID", user.ID.String()),
			zap.String("phone", user.Phone),
		)

		c.Locals("user", user)
		c.Locals("session", session)
		c.Locals("userID", user.ID)

		return c.Next()
	}
}

func GetUserFromContext(c *fiber.Ctx) (*port.User, bool) {
	user, ok := c.Locals("user").(*port.User)
	return user, ok
}

func GetUserIDFromContext(c *fiber.Ctx) (uuid.UUID, bool) {
	userID, ok := c.Locals("userID").(uuid.UUID)
	return userID, ok
}

func GetSessionFromContext(c *fiber.Ctx) (*auth.Session, bool) {
	session, ok := c.Locals("session").(*auth.Session)
	return session, ok
}

func RequireAuth(c *fiber.Ctx) (*port.User, bool) {
	user, ok := GetUserFromContext(c)
	if !ok {
		c.Status(
			fiber.StatusUnauthorized,
		).JSON(fiber.Map{"error": "Authentication required"})
		return nil, false
	}
	return user, true
}

func RequirePhoneVerification(c *fiber.Ctx) (*port.User, bool) {
	user, ok := RequireAuth(c)
	if !ok {
		return nil, false
	}

	if user.Verified < 1 {
		c.Status(
			fiber.StatusForbidden,
		).JSON(fiber.Map{"error": "Phone verification required"})
		return nil, false
	}

	return user, true
}
