package port

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PaymentService interface {
	InitiatePayment(ctx context.Context, req *InitiatePaymentRequest) (*InitiatePaymentResponse, error)
	VerifyPayment(ctx context.Context, paymentID uuid.UUID) (*VerifyPaymentResponse, error)
}

type InitiatePaymentRequest struct {
	OrderID     uuid.UUID              `json:"order_id"`
	UserID      uuid.UUID              `json:"user_id"`
	Amount      int                    `json:"amount"`
	CallbackURL string                 `json:"callback_url"`
	ReturnURL   string                 `json:"return_url"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type InitiatePaymentResponse struct {
	PaymentID    uuid.UUID `json:"payment_id"`
	GatewayURL   string    `json:"gateway_url"`
	GatewayRefID string    `json:"gateway_ref_id"`
	Status       string    `json:"status"`
	Amount       int       `json:"amount"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type VerifyPaymentResponse struct {
	PaymentID    uuid.UUID  `json:"payment_id"`
	Status       string     `json:"status"`
	Amount       int        `json:"amount"`
	GatewayRefID string     `json:"gateway_ref_id,omitempty"`
	PaidAt       *time.Time `json:"paid_at,omitempty"`
	FailedAt     *time.Time `json:"failed_at,omitempty"`
}
