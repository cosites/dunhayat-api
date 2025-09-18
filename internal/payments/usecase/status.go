package usecase

import (
	"context"
	"errors"
	"fmt"

	"dunhayat-api/internal/payments"
	"dunhayat-api/internal/payments/port"

	"github.com/google/uuid"
)

type GetPaymentStatusUseCase interface {
	Execute(
		ctx context.Context,
		req *payments.GetPaymentStatusRequest,
	) (*payments.GetPaymentStatusResponse, error)
}

type getPaymentStatusUseCase struct {
	orderService port.OrderService
}

func NewGetPaymentStatusUseCase(
	orderService port.OrderService,
) GetPaymentStatusUseCase {
	return &getPaymentStatusUseCase{
		orderService: orderService,
	}
}

func (uc *getPaymentStatusUseCase) Execute(
	ctx context.Context,
	req *payments.GetPaymentStatusRequest,
) (*payments.GetPaymentStatusResponse, error) {
	var sale *port.Sale
	var err error

	if req.TrackingCode != "" {
		sale, err = uc.orderService.GetSaleByTrackingCode(
			ctx, req.TrackingCode,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get sale by tracking code: %w", err,
			)
		}
	} else if req.OrderID != "" {
		orderID, parseErr := uuid.Parse(req.OrderID)
		if parseErr != nil {
			return nil, fmt.Errorf(
				"invalid order ID format: %w", parseErr,
			)
		}
		sale, err = uc.orderService.GetSaleByID(ctx, orderID)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get sale by ID: %w", err,
			)
		}
	} else {
		return nil, errors.New(
			"either tracking code or order ID must be provided",
		)
	}

	if sale == nil {
		return nil, errors.New(
			"sale not found",
		)
	}

	var paymentStatus payments.PaymentStatus
	switch sale.Status {
	case port.OrderStatusPending:
		paymentStatus = payments.PaymentStatusPending
	case port.OrderStatusPaid:
		paymentStatus = payments.PaymentStatusPaid
	case port.OrderStatusFailed:
		paymentStatus = payments.PaymentStatusFailed
	case port.OrderStatusCancelled:
		paymentStatus = payments.PaymentStatusCancelled
	default:
		paymentStatus = payments.PaymentStatusPending
	}

	return &payments.GetPaymentStatusResponse{
		PaymentID:    sale.ID,
		TrackingCode: sale.TrackingCode,
		Status:       paymentStatus,
		Amount:       sale.TotalPrice,
		OrderStatus:  string(sale.Status),
		CreatedAt:    sale.CreatedAt,
		UpdatedAt:    sale.UpdatedAt,
	}, nil
}
