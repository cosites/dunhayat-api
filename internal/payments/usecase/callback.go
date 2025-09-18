package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"dunhayat-api/internal/payments"
	"dunhayat-api/internal/payments/port"
)

type HandleCallbackUseCase interface {
	Execute(
		ctx context.Context,
		callbackData payments.PaymentCallbackRequest,
	) error
}

type handleCallbackUseCase struct {
	orderService port.OrderService
}

func NewHandleCallbackUseCase(
	orderService port.OrderService,
) HandleCallbackUseCase {
	return &handleCallbackUseCase{
		orderService: orderService,
	}
}

func (uc *handleCallbackUseCase) Execute(
	ctx context.Context,
	callbackData payments.PaymentCallbackRequest,
) error {
	sale, err := uc.orderService.GetSaleByTrackingCode(
		ctx, callbackData.TrackID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to get sale by tracking code: %w", err,
		)
	}
	if sale == nil {
		return errors.New(
			"sale not found for track ID: " + callbackData.TrackID,
		)
	}

	// XXX: Reserved for future audit store
	_, _ = json.Marshal(callbackData)

	var newStatus port.OrderStatus
	if callbackData.Success && callbackData.Status == 100 {
		newStatus = port.OrderStatusPaid
	} else {
		newStatus = port.OrderStatusFailed
	}

	if err := uc.orderService.UpdateSaleStatus(
		ctx, sale.ID, newStatus,
	); err != nil {
		return fmt.Errorf(
			"failed to update sale status: %w", err,
		)
	}

	// XXX: Reserved for future timestamps
	_ = time.Now()

	return nil
}
