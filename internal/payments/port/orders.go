package port

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusFailed    OrderStatus = "failed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

func (s OrderStatus) String() string {
	return string(s)
}

type Sale struct {
	ID           uuid.UUID   `json:"id"`
	UserID       uuid.UUID   `json:"user_id"`
	Status       OrderStatus `json:"status"`
	TrackingCode *string     `json:"tracking_code,omitempty"`
	TotalPrice   int         `json:"total_price"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type OrderPort interface {
	GetSaleByID(ctx context.Context, saleID uuid.UUID) (*Sale, error)
	GetSaleByTrackingCode(ctx context.Context, trackingCode string) (*Sale, error)
	UpdateSaleStatus(ctx context.Context, saleID uuid.UUID, status OrderStatus) error
	SetSaleTrackingCode(ctx context.Context, saleID uuid.UUID, trackingCode string) error
}
