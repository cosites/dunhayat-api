package adapter

import (
	"context"

	"dunhayat-api/internal/orders"
	"dunhayat-api/internal/orders/repository"
	"dunhayat-api/internal/payments/port"

	"github.com/google/uuid"
)

type PaymentsOrderService struct {
	saleRepo repository.SaleRepository
}

func NewPaymentsOrderService(
	saleRepo repository.SaleRepository,
) port.OrderService {
	return &PaymentsOrderService{
		saleRepo: saleRepo,
	}
}

func (s *PaymentsOrderService) GetSaleByID(
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

func (s *PaymentsOrderService) GetSaleByTrackingCode(
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

func (s *PaymentsOrderService) UpdateSaleStatus(
	ctx context.Context,
	saleID uuid.UUID,
	status port.OrderStatus,
) error {
	return s.saleRepo.UpdateStatus(ctx, saleID, orders.OrderStatus(status))
}

func (s *PaymentsOrderService) SetSaleTrackingCode(
	ctx context.Context,
	saleID uuid.UUID,
	trackingCode string,
) error {
	return s.saleRepo.SetTrackingCode(ctx, saleID, trackingCode)
}
