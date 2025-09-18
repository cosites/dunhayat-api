package router

import (
	"context"
	"strings"
	"time"

	"dunhayat-api/pkg/config"
	"dunhayat-api/pkg/logger"
	"dunhayat-api/pkg/router/port"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"
)

type FiberRouter struct {
	app            *fiber.App
	logger         logger.Interface
	cfg            *FiberConfig
	productHandler port.ProductHandler
	orderHandler   port.OrderHandler
	paymentHandler port.PaymentHandler
	authHandler    port.AuthHandler
	authMiddleware port.AuthMiddleware
	version        string
}

type FiberConfig struct {
	AppEnv string
	CORS   *config.CORSConfig
}

func NewFiberRouter(
	log logger.Interface,
	cfg *FiberConfig,
	productHandler port.ProductHandler,
	orderHandler port.OrderHandler,
	paymentHandler port.PaymentHandler,
	authHandler port.AuthHandler,
	authMiddleware port.AuthMiddleware,
	version string,
) *FiberRouter {
	app := fiber.New(fiber.Config{
		AppName:                 "Dunhayat API",
		EnableTrustedProxyCheck: true,
		ProxyHeader:             "X-Forwarded-For",
		ReadTimeout:             30 * time.Second,
		WriteTimeout:            30 * time.Second,
		IdleTimeout:             120 * time.Second,
		ReadBufferSize:          8192,
		WriteBufferSize:         8192,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			log.Error("Unhandled HTTP error",
				zap.Error(err),
				zap.Int("status_code", code),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.String("ip", c.IP()),
				zap.String("user_agent", c.Get("User-Agent")),
			)

			return c.Status(code).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
				"code":    code,
			})
		},
		DisableStartupMessage: true,
	})

	router := &FiberRouter{
		app:            app,
		logger:         log,
		cfg:            cfg,
		productHandler: productHandler,
		orderHandler:   orderHandler,
		paymentHandler: paymentHandler,
		authHandler:    authHandler,
		authMiddleware: authMiddleware,
		version:        version,
	}

	router.setupMiddleware()
	router.setupRoutes()

	return router
}

func (r *FiberRouter) setupMiddleware() {
	r.app.Use(recover.New(recover.Config{
		EnableStackTrace: r.cfg.AppEnv == "development",
	}))

	if r.cfg.CORS != nil {
		r.app.Use(cors.New(cors.Config{
			AllowOrigins:     strings.Join(r.cfg.CORS.AllowedOrigins, ","),
			AllowMethods:     strings.Join(r.cfg.CORS.AllowedMethods, ","),
			AllowHeaders:     strings.Join(r.cfg.CORS.AllowedHeaders, ","),
			AllowCredentials: r.cfg.CORS.AllowCredentials,
			MaxAge:           300, // 5 minutes
		}))
	}

	r.app.Use(func(c *fiber.Ctx) error {
		start := time.Now()

		if c.Path() == "/health" || c.Path() == "/" {
			return c.Next()
		}

		r.logger.Info("HTTP Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("query", c.Query("")),
			zap.String("remote_addr", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
			zap.String("referer", c.Get("Referer")),
			zap.Int(
				"content_length",
				c.Request().Header.ContentLength(),
			),
		)

		err := c.Next()

		duration := time.Since(start)
		statusCode := c.Response().StatusCode()

		fields := []zap.Field{
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status_code", statusCode),
			zap.Duration("duration", duration),
		}

		switch {
		case statusCode >= 500:
			r.logger.Error("HTTP Request Error", fields...)
		case statusCode >= 400:
			r.logger.Warn("HTTP Request Warning", fields...)
		case statusCode >= 300:
			r.logger.Info("HTTP Request Redirect", fields...)
		default:
			r.logger.Info("HTTP Request", fields...)
		}

		return err
	})
}

func (r *FiberRouter) setupRoutes() {
	r.app.Get("/health", r.handleHealth)

	r.app.Get("/version", r.handleVersion)

	if r.cfg.AppEnv == "development" {
		r.logger.Info(
			"Development mode detected - Swagger UI enabled",
		)
		r.setupSwaggerRoutes()
	}

	api := r.app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post(
		"/request-otp",
		r.authHandler.RequestOTP,
	)
	auth.Post(
		"/verify-otp",
		r.authHandler.VerifyOTP,
	)
	auth.Get(
		"/otp-status",
		r.authHandler.GetOTPStatus,
	)
	auth.Post(
		"/logout",
		r.authHandler.Logout,
	)

	products := api.Group("/products")
	products.Get(
		"/",
		r.productHandler.ListProducts,
	)
	products.Get(
		"/:id",
		r.productHandler.GetProduct,
	)

	orders := api.Group("/orders")
	orders.Post(
		"/",
		r.authMiddleware.Authenticate(),
		r.orderHandler.CreateOrder,
	)
	orders.Get(
		"/:id",
		r.authMiddleware.Authenticate(),
		r.orderHandler.GetOrder,
	)

	payments := api.Group("/payments")
	payments.Post(
		"/initiate",
		r.authMiddleware.Authenticate(),
		r.paymentHandler.InitiatePayment,
	)
	payments.Post(
		"/verify",
		r.authMiddleware.Authenticate(),
		r.paymentHandler.VerifyPayment,
	)
	payments.Post(
		"/callback",
		r.paymentHandler.HandleCallback,
	)
	payments.Get(
		"/:id/status",
		r.authMiddleware.Authenticate(),
		r.paymentHandler.GetPaymentStatus,
	)
}

func (r *FiberRouter) setupSwaggerRoutes() {
	r.app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Swagger JSON endpoint - implementation pending",
		})
	})

	r.app.Get("/swagger/*", swagger.HandlerDefault)
}

func (r *FiberRouter) handleHealth(c *fiber.Ctx) error {
	r.logger.Info("Health check requested",
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.String("status", "200 OK"),
	)

	return c.JSON(fiber.Map{
		"status":    "success",
		"message":   "Service is healthy",
		"timestamp": time.Now().UTC(),
		"version":   r.version,
	})
}

func (r *FiberRouter) handleVersion(c *fiber.Ctx) error {
	r.logger.Info("Version info requested",
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.String("version", r.version),
	)

	return c.JSON(fiber.Map{
		"status":    "success",
		"message":   "Version information",
		"version":   r.version,
		"timestamp": time.Now().UTC(),
	})
}

func (r *FiberRouter) GetApp() *fiber.App {
	return r.app
}

func (r *FiberRouter) Start(addr string) error {
	return r.app.Listen(addr)
}

func (r *FiberRouter) Shutdown(ctx context.Context) error {
	return r.app.ShutdownWithContext(ctx)
}
