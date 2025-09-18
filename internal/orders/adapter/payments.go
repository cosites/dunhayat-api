package adapter

import (
	"context"

	"dunhayat-api/internal/orders"
	"dunhayat-api/internal/orders/repository"
	"dunhayat-api/internal/payments/port"

	"github.com/google/uuid"
)

type PaymentsOrderAdapter struct {
	saleRepo repository.SaleRepository
}

func NewPaymentsOrderAdapter(
	saleRepo repository.SaleRepository,
) port.OrderPort {
	return &PaymentsOrderAdapter{
		saleRepo: saleRepo,
	}
}

func (s *PaymentsOrderAdapter) GetSaleByID(
	ctx context.Context,
	saleID uuid.UUID,
) (*port.Sale, error) {
	sale, err := s.saleRepo.GetByID(ctx, saleID)
	if err != nil {
		return nil, err
	}
	if sale == nil {
		return nil, nil
	}

	return &port.Sale{
		ID:           sale.ID,
		UserID:       sale.UserID,
		Status:       port.OrderStatus(sale.Status),
		TrackingCode: sale.TrackingCode,
		TotalPrice:   sale.TotalPrice,
		CreatedAt:    sale.CreatedAt,
		UpdatedAt:    sale.UpdatedAt,
	}, nil
}

func (s *PaymentsOrderAdapter) GetSaleByTrackingCode(
	ctx context.Context,
	trackingCode string,
) (*port.Sale, error) {
	sale, err := s.saleRepo.GetByTrackingCode(ctx, trackingCode)
	if err != nil {
		return nil, err
	}
	if sale == nil {
		return nil, nil
	}

	return &port.Sale{
		ID:           sale.ID,
		UserID:       sale.UserID,
		Status:       port.OrderStatus(sale.Status),
		TrackingCode: sale.TrackingCode,
		TotalPrice:   sale.TotalPrice,
		CreatedAt:    sale.CreatedAt,
		UpdatedAt:    sale.UpdatedAt,
	}, nil
}

func (s *PaymentsOrderAdapter) UpdateSaleStatus(
	ctx context.Context,
	saleID uuid.UUID,
	status port.OrderStatus,
) error {
	return s.saleRepo.UpdateStatus(ctx, saleID, orders.OrderStatus(status))
}

func (s *PaymentsOrderAdapter) SetSaleTrackingCode(
	ctx context.Context,
	saleID uuid.UUID,
	trackingCode string,
) error {
	return s.saleRepo.SetTrackingCode(ctx, saleID, trackingCode)
}
