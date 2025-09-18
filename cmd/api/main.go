package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	authHandler "dunhayat-api/internal/auth/http"
	authLogger "dunhayat-api/internal/auth/logger"
	authRepo "dunhayat-api/internal/auth/repository"
	authUseCase "dunhayat-api/internal/auth/usecase"
	orderAdapter "dunhayat-api/internal/orders/adapter"
	orderHandler "dunhayat-api/internal/orders/http"
	orderRepo "dunhayat-api/internal/orders/repository"
	orderUseCase "dunhayat-api/internal/orders/usecase"
	paymentAdapter "dunhayat-api/internal/payments/adapter"
	paymentHandler "dunhayat-api/internal/payments/http"
	paymentUseCase "dunhayat-api/internal/payments/usecase"
	productAdapter "dunhayat-api/internal/products/adapter"
	productHandler "dunhayat-api/internal/products/http"
	productRepo "dunhayat-api/internal/products/repository"
	productUseCase "dunhayat-api/internal/products/usecase"
	usersAdapter "dunhayat-api/internal/users/adapter"
	userRepo "dunhayat-api/internal/users/repository"
	"dunhayat-api/pkg/config"
	"dunhayat-api/pkg/database"
	"dunhayat-api/pkg/logger"
	"dunhayat-api/pkg/payment"
	"dunhayat-api/pkg/redis"
	"dunhayat-api/pkg/router"
	"dunhayat-api/pkg/sms"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var Version = "dev"

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Failed to load configuration: %v\n",
			err,
		)
		os.Exit(1)
	}

	var env logger.Env
	if cfg.Env == "production" {
		env = logger.EnvProduction
	} else {
		env = logger.EnvDevelopment
	}

	log := logger.New(env, cfg.Log.Level, uuid.New())
	log.Info("Starting Dunhayat Coffee Roastery API...")
	log.Info("Configuration loaded successfully")

	dbConn, err := database.Connect(&cfg.Database, log)
	if err != nil {
		log.Fatal(
			"Failed to connect to database",
			zap.Error(err),
		)
	}
	defer func() {
		if sqlDB, err := dbConn.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Error(
					"Failed to close database connection",
					zap.Error(err),
				)
			}
			log.Info("Database connection closed successfully")
		} else {
			log.Error(
				"Failed to retrieve underlying SQL DB from dbConn",
				zap.Error(err),
			)
		}
	}()

	redisClient, err := redis.Connect(&redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}, log)
	if err != nil {
		log.Fatal(
			"Failed to connect to Redis",
			zap.Error(err),
		)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Error(
				"Failed to close Redis connection",
				zap.Error(err),
			)
			return
		}
		log.Info("Redis connection closed successfully")
	}()

	log.Info(
		"Database connection established - migrations handled by Atlas",
	)

	userRepository := userRepo.NewUserRepository(
		dbConn,
	)
	addressRepository := userRepo.NewAddressRepository(
		dbConn,
	)
	productRepository := productRepo.NewProductRepository(
		dbConn,
	)
	saleRepository := orderRepo.NewSaleRepository(
		dbConn,
	)
	saleItemRepository := orderRepo.NewSaleItemRepository(
		dbConn,
	)
	cartReservationRepository := orderRepo.NewCartReservationRepository(
		dbConn,
	)
	sessionRepository := authRepo.NewSessionRepository(
		dbConn,
	)

	smsProvider := sms.NewKavenegarProvider(
		cfg.Auth.KavenegarAPIKey,
	)
	authLoggerAdapter := authLogger.NewAdapter(
		log,
	)
	otpRepository := authRepo.NewRedisOTPRepository(
		redisClient, log,
	)

	ordersProductService := productAdapter.NewOrdersProductService(
		productRepository,
	)
	ordersUserService := usersAdapter.NewOrdersUserService(
		userRepository,
	)

	listProductsUseCase := productUseCase.NewListProductsUseCase(
		productRepository,
	)
	getProductUseCase := productUseCase.NewGetProductUseCase(
		productRepository,
	)

	zibalClient := payment.NewZibalClient(payment.ZibalConfig{
		MerchantID: cfg.Payment.Zibal.MerchantID,
		BaseURL:    cfg.Payment.Zibal.BaseURL,
		Timeout:    time.Duration(cfg.Payment.Zibal.Timeout) * time.Second,
	})

	paymentsOrderService := orderAdapter.NewPaymentsOrderService(
		saleRepository,
	)

	initiatePaymentUseCase := paymentUseCase.NewInitiatePaymentUseCase(
		paymentsOrderService,
		zibalClient,
	)
	verifyPaymentUseCase := paymentUseCase.NewVerifyPaymentUseCase(
		paymentsOrderService,
		zibalClient,
	)
	handleCallbackUseCase := paymentUseCase.NewHandleCallbackUseCase(
		paymentsOrderService,
	)
	getPaymentStatusUseCase := paymentUseCase.NewGetPaymentStatusUseCase(
		paymentsOrderService,
	)

	createOrderUseCase := orderUseCase.NewCreateOrderUseCase(
		saleRepository,
		saleItemRepository,
		cartReservationRepository,
		ordersProductService,
		ordersUserService,
		paymentAdapter.NewOrdersPaymentService(
			initiatePaymentUseCase,
			verifyPaymentUseCase,
		),
	)

	authUserService := usersAdapter.NewAuthUserService(
		userRepository, addressRepository,
	)

	requestOTPUseCase := authUseCase.NewRequestOTPUseCase(
		otpRepository,
		smsProvider,
		cfg.Auth.OTPTemplate,
		authLoggerAdapter,
	)
	verifyOTPUseCase := authUseCase.NewVerifyOTPUseCase(
		otpRepository,
		sessionRepository,
		authUserService,
		24*time.Hour, // XXX: could be configurable
		authLoggerAdapter,
	)
	logoutUseCase := authUseCase.NewLogoutUseCase(
		sessionRepository,
		authLoggerAdapter,
	)

	productHTTPHandler := productHandler.NewProductHandler(
		listProductsUseCase,
		getProductUseCase,
	)
	orderHTTPHandler := orderHandler.NewOrderHandler(
		createOrderUseCase,
	)
	paymentHTTPHandler := paymentHandler.NewPaymentHandler(
		initiatePaymentUseCase,
		verifyPaymentUseCase,
		handleCallbackUseCase,
		getPaymentStatusUseCase,
	)
	authHTTPHandler := authHandler.NewAuthHandler(
		requestOTPUseCase,
		verifyOTPUseCase,
		logoutUseCase,
	)

	middlewareUserService := usersAdapter.NewMiddlewareUserService(
		userRepository,
	)
	authMiddleware := authHandler.NewAuthMiddleware(
		log,
		sessionRepository,
		middlewareUserService,
	)

	version := Version

	log.Info("Application version", zap.String("version", version))

	routerConfig := &router.FiberConfig{
		AppEnv: cfg.Env,
		CORS:   &cfg.CORS,
	}
	fiberRouter := router.NewFiberRouter(
		log,
		routerConfig,
		productHTTPHandler,
		orderHTTPHandler,
		paymentHTTPHandler,
		authHTTPHandler,
		authMiddleware,
		version,
	)

	serverAddr := fmt.Sprintf(
		"%s:%s", cfg.Server.Host, cfg.Server.Port,
	)
	log.Info(
		"Fiber server configured",
		zap.String("host", cfg.Server.Host),
		zap.String("port", cfg.Server.Port),
		zap.String("environment", cfg.Env),
		zap.String("logLevel", cfg.Log.Level),
	)

	go func() {
		log.Info(
			"Starting Fiber server",
			zap.String("host", cfg.Server.Host),
			zap.String("port", cfg.Server.Port),
			zap.String("environment", cfg.Env),
			zap.String("logLevel", cfg.Log.Level),
		)
		if err := fiberRouter.Start(serverAddr); err != nil {
			log.Fatal(
				"Failed to start server", zap.Error(err),
			)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(), 30*time.Second,
	)
	defer cancel()

	if err := fiberRouter.Shutdown(ctx); err != nil {
		log.Fatal(
			"Server forced to shutdown", zap.Error(err),
		)
	}

	log.Info("Server shutdown completed successfully")
}
