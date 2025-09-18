package payments

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

func (s PaymentStatus) String() string {
	return string(s)
}

type PaymentMethod string

const (
	PaymentMethodZibal PaymentMethod = "zibal"
)

func (m PaymentMethod) String() string {
	return string(m)
}

type Payment struct {
	ID           uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID      uuid.UUID     `json:"order_id" gorm:"type:uuid;not null"`
	UserID       uuid.UUID     `json:"user_id" gorm:"type:uuid;not null"`
	Amount       int           `json:"amount" gorm:"not null;check:amount > 0"`
	Status       PaymentStatus `json:"status" gorm:"type:varchar(50);not null;default:'pending'"`
	Method       PaymentMethod `json:"method" gorm:"type:varchar(50);not null"`
	GatewayRefID *string       `json:"gateway_ref_id,omitempty" gorm:"type:varchar(255)"`
	GatewayURL   *string       `json:"gateway_url,omitempty" gorm:"type:text"`
	CallbackURL  string        `json:"callback_url" gorm:"type:text;not null"`
	ReturnURL    string        `json:"return_url" gorm:"type:text;not null"`
	Description  string        `json:"description" gorm:"type:text"`
	Metadata     *string       `json:"metadata,omitempty" gorm:"type:jsonb"`
	PaidAt       *time.Time    `json:"paid_at,omitempty"`
	FailedAt     *time.Time    `json:"failed_at,omitempty"`
	CreatedAt    time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

type PaymentCallback struct {
	ID          uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PaymentID   uuid.UUID     `json:"payment_id" gorm:"type:uuid;not null"`
	GatewayData string        `json:"gateway_data" gorm:"type:jsonb;not null"`
	Status      PaymentStatus `json:"status" gorm:"type:varchar(50);not null"`
	ProcessedAt time.Time     `json:"processed_at" gorm:"autoCreateTime"`
	CreatedAt   time.Time     `json:"created_at" gorm:"autoCreateTime"`
}

type InitiatePaymentRequest struct {
	OrderID     uuid.UUID      `json:"order_id" binding:"required"`
	UserID      uuid.UUID      `json:"user_id" binding:"required"`
	Amount      int            `json:"amount" binding:"required,min=1"`
	Method      PaymentMethod  `json:"method" binding:"required"`
	CallbackURL string         `json:"callback_url" binding:"required"`
	ReturnURL   string         `json:"return_url" binding:"required"`
	Description string         `json:"description"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

type InitiatePaymentResponse struct {
	PaymentID    uuid.UUID     `json:"payment_id"`
	GatewayURL   string        `json:"gateway_url"`
	GatewayRefID string        `json:"gateway_ref_id"`
	Status       PaymentStatus `json:"status"`
	Amount       int           `json:"amount"`
	ExpiresAt    time.Time     `json:"expires_at"`
}

type VerifyPaymentRequest struct {
	PaymentID uuid.UUID `json:"payment_id" binding:"required"`
}

type VerifyPaymentResponse struct {
	PaymentID    uuid.UUID     `json:"payment_id"`
	Status       PaymentStatus `json:"status"`
	Amount       int           `json:"amount"`
	GatewayRefID string        `json:"gateway_ref_id,omitempty"`
	PaidAt       *time.Time    `json:"paid_at,omitempty"`
	FailedAt     *time.Time    `json:"failed_at,omitempty"`
}

type PaymentCallbackRequest struct {
	Success          bool   `json:"success"`
	Status           int    `json:"status"`
	TrackID          string `json:"trackId"`
	OrderID          string `json:"orderId"`
	Amount           int    `json:"amount"`
	CardNumber       string `json:"cardNumber"`
	HashedCardNumber string `json:"hashedCardNumber"`
}

type GetPaymentStatusRequest struct {
	OrderID      string `json:"order_id,omitempty"`
	TrackingCode string `json:"tracking_code,omitempty"`
}

type GetPaymentStatusResponse struct {
	PaymentID    uuid.UUID     `json:"payment_id"`
	TrackingCode *string       `json:"tracking_code,omitempty"`
	Status       PaymentStatus `json:"status"`
	Amount       int           `json:"amount"`
	OrderStatus  string        `json:"order_status"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

func (Payment) TableName() string {
	return "payments"
}

func (PaymentCallback) TableName() string {
	return "payment_callbacks"
}
