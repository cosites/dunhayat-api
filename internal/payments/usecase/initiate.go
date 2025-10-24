package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"dunhayat-api/internal/payments"
	"dunhayat-api/internal/payments/port"
	"dunhayat-api/pkg/config"
	"dunhayat-api/pkg/logger"
	"dunhayat-api/pkg/payment"

	"go.uber.org/zap"
)

type InitiatePaymentUseCase interface {
	Execute(
		ctx context.Context,
		req *payments.InitiatePaymentRequest,
	) (*payments.InitiatePaymentResponse, error)
}

type initiatePaymentUseCase struct {
	orderPort   port.OrderPort
	zibalClient *payment.ZibalClient
	logger      logger.Interface
	config      *config.Config
}

func NewInitiatePaymentUseCase(
	orderPort port.OrderPort,
	zibalClient *payment.ZibalClient,
	logger logger.Interface,
	config *config.Config,
) InitiatePaymentUseCase {
	return &initiatePaymentUseCase{
		orderPort:   orderPort,
		zibalClient: zibalClient,
		logger:      logger,
		config:      config,
	}
}

func (uc *initiatePaymentUseCase) Execute(
	ctx context.Context,
	req *payments.InitiatePaymentRequest,
) (*payments.InitiatePaymentResponse, error) {
	callbackURL := uc.config.App.Domain + req.CallbackURL

	zibalReq := payment.ZibalPaymentRequest{
		Amount:      req.Amount,
		OrderID:     req.OrderID.String(),
		CallbackURL: callbackURL,
		Description: req.Description,
	}

	uc.logger.Info("Initiating Zibal payment request",
		zap.String("order_id", req.OrderID.String()),
		zap.Int("amount", req.Amount),
		zap.String("callback_url", callbackURL),
	)

	zibalResp, err := uc.zibalClient.CreatePaymentRequest(zibalReq)
	if err != nil {
		uc.logger.Error("Zibal payment request failed",
			zap.String("order_id", req.OrderID.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf(
			"failed to create zibal payment request: %w", err,
		)
	}

	trackIDStr := strconv.FormatInt(zibalResp.TrackID, 10)
	uc.logger.Info("Zibal payment request successful",
		zap.String("order_id", req.OrderID.String()),
		zap.String("track_id", trackIDStr),
	)

	if err := uc.orderPort.SetSaleTrackingCode(
		ctx, req.OrderID, trackIDStr,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to set tracking code on sale: %w", err,
		)
	}

	gatewayURL := uc.zibalClient.GetPaymentURL(zibalResp.TrackID)

	return &payments.InitiatePaymentResponse{
		PaymentID:    req.OrderID,
		GatewayURL:   gatewayURL,
		GatewayRefID: trackIDStr,
		Status:       payments.PaymentStatusPending,
		Amount:       req.Amount,
		ExpiresAt:    time.Now().Add(30 * time.Minute),
	}, nil
}
