package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"dunhayat-api/internal/payments"
	"dunhayat-api/internal/payments/port"
	"dunhayat-api/pkg/payment"
)

type VerifyPaymentUseCase interface {
	Execute(
		ctx context.Context,
		req *payments.VerifyPaymentRequest,
	) (*payments.VerifyPaymentResponse, error)
}

type verifyPaymentUseCase struct {
	orderService port.OrderService
	zibalClient  *payment.ZibalClient
}

func NewVerifyPaymentUseCase(
	orderService port.OrderService,
	zibalClient *payment.ZibalClient,
) VerifyPaymentUseCase {
	return &verifyPaymentUseCase{
		orderService: orderService,
		zibalClient:  zibalClient,
	}
}

func (uc *verifyPaymentUseCase) Execute(
	ctx context.Context,
	req *payments.VerifyPaymentRequest,
) (*payments.VerifyPaymentResponse, error) {
	sale, err := uc.orderService.GetSaleByID(ctx, req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sale: %w", err)
	}
	if sale == nil {
		return nil, errors.New("sale not found")
	}
	if sale.TrackingCode == nil || *sale.TrackingCode == "" {
		return nil, errors.New("sale has no tracking code")
	}

	zibalReq := payment.ZibalVerifyRequest{TrackID: *sale.TrackingCode}
	zibalResp, err := uc.zibalClient.VerifyPayment(zibalReq)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to verify payment with zibal: %w", err,
		)
	}

	var newStatus port.OrderStatus
	switch zibalResp.Result {
	case 100:
		newStatus = port.OrderStatusPaid
	case 102:
		newStatus = port.OrderStatusFailed
	default:
		newStatus = port.OrderStatusFailed
	}

	if err := uc.orderService.UpdateSaleStatus(
		ctx, sale.ID, newStatus,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to update sale status: %w", err,
		)
	}

	var paidAt *time.Time
	var failedAt *time.Time
	now := time.Now()

	switch newStatus {
	case port.OrderStatusPaid:
		paidAt = &now
	case port.OrderStatusFailed:
		failedAt = &now
	}

	return &payments.VerifyPaymentResponse{
		PaymentID:    sale.ID,
		Status:       payments.PaymentStatus(newStatus),
		Amount:       sale.TotalPrice,
		GatewayRefID: *sale.TrackingCode,
		PaidAt:       paidAt,
		FailedAt:     failedAt,
	}, nil
}
