package usecase

import (
	"context"
	"fmt"
	"time"

	"dunhayat-api/internal/payments"
	"dunhayat-api/internal/payments/port"
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
}

func NewInitiatePaymentUseCase(
	orderPort port.OrderPort,
	zibalClient *payment.ZibalClient,
	logger logger.Interface,
) InitiatePaymentUseCase {
	return &initiatePaymentUseCase{
		orderPort:   orderPort,
		zibalClient: zibalClient,
		logger:      logger,
	}
}

func (uc *initiatePaymentUseCase) Execute(
	ctx context.Context,
	req *payments.InitiatePaymentRequest,
) (*payments.InitiatePaymentResponse, error) {
	zibalReq := payment.ZibalPaymentRequest{
		Amount:      req.Amount,
		OrderID:     req.OrderID.String(),
		CallbackURL: req.CallbackURL,
		Description: req.Description,
	}

	uc.logger.Info("Initiating Zibal payment request",
		zap.String("order_id", req.OrderID.String()),
		zap.Int("amount", req.Amount),
		zap.String("callback_url", req.CallbackURL),
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

	uc.logger.Info("Zibal payment request successful",
		zap.String("order_id", req.OrderID.String()),
		zap.String("track_id", zibalResp.TrackID),
	)

	if err := uc.orderPort.SetSaleTrackingCode(
		ctx, req.OrderID, zibalResp.TrackID,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to set tracking code on sale: %w", err,
		)
	}

	gatewayURL := uc.zibalClient.GetPaymentURL(zibalResp.TrackID)

	return &payments.InitiatePaymentResponse{
		PaymentID:    req.OrderID,
		GatewayURL:   gatewayURL,
		GatewayRefID: zibalResp.TrackID,
		Status:       payments.PaymentStatusPending,
		Amount:       req.Amount,
		ExpiresAt:    time.Now().Add(30 * time.Minute),
	}, nil
}
