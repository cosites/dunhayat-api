package usecase

import (
	"context"
	"fmt"
	"time"

	"dunhayat-api/internal/payments"
	"dunhayat-api/internal/payments/port"
	"dunhayat-api/pkg/payment"
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
}

func NewInitiatePaymentUseCase(
	orderPort port.OrderPort,
	zibalClient *payment.ZibalClient,
) InitiatePaymentUseCase {
	return &initiatePaymentUseCase{
		orderPort:   orderPort,
		zibalClient: zibalClient,
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

	zibalResp, err := uc.zibalClient.CreatePaymentRequest(zibalReq)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create zibal payment request: %w", err,
		)
	}

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
