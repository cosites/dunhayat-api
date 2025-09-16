package orders

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusFailed    OrderStatus = "failed"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
)

func (s OrderStatus) String() string {
	return string(s)
}

type Sale struct {
	ID           uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID   `json:"user_id" gorm:"type:uuid;not null"`
	Status       OrderStatus `json:"status" gorm:"type:varchar(50);not null;default:'pending'"`
	TrackingCode *string     `json:"tracking_code,omitempty"`
	TotalPrice   int         `json:"total_price" gorm:"not null;check:total_price > 0"`
	CreatedAt    time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

type SaleItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SaleID    uuid.UUID `json:"sale_id" gorm:"type:uuid;not null"`
	ProductID string    `json:"product_id" gorm:"type:varchar(100);not null"`
	Quantity  int       `json:"quantity" gorm:"not null;check:quantity > 0"`
	Price     int       `json:"price" gorm:"not null;check:price > 0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type CartReservation struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ProductID string    `json:"product_id" gorm:"type:varchar(100);not null"`
	Quantity  int       `json:"quantity" gorm:"not null;check:quantity > 0"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type CreateOrderRequest struct {
	UserID      uuid.UUID          `json:"user_id" binding:"required"`
	Items       []OrderItemRequest `json:"items" binding:"required,min=1"`
	Address     string             `json:"address" binding:"required"`
	PostalCode  string             `json:"postal_code" binding:"required"`
	CallbackURL string             `json:"callback_url" binding:"required"`
	ReturnURL   string             `json:"return_url" binding:"required"`
}

type OrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type OrderResponse struct {
	ID           uuid.UUID    `json:"id"`
	UserID       uuid.UUID    `json:"user_id"`
	Status       OrderStatus  `json:"status"`
	TrackingCode *string      `json:"tracking_code,omitempty"`
	TotalPrice   int          `json:"total_price"`
	Items        []SaleItem   `json:"items"`
	Address      string       `json:"address"`
	PostalCode   string       `json:"postal_code"`
	Payment      *PaymentInfo `json:"payment,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type PaymentInfo struct {
	PaymentID    uuid.UUID `json:"payment_id"`
	GatewayURL   string    `json:"gateway_url"`
	GatewayRefID string    `json:"gateway_ref_id"`
	Status       string    `json:"status"`
	Amount       int       `json:"amount"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (Sale) TableName() string {
	return "sales"
}

func (SaleItem) TableName() string {
	return "sale_items"
}

func (CartReservation) TableName() string {
	return "cart_reservations"
}
