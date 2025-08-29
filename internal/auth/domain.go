package auth

import (
	"time"

	"github.com/google/uuid"
)

type OTPStatus string

const (
	OTPStatusPending  OTPStatus = "pending"
	OTPStatusVerified OTPStatus = "verified"
	OTPStatusExpired  OTPStatus = "expired"
	OTPStatusFailed   OTPStatus = "failed"
)

type OTP struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Phone     string    `json:"phone" gorm:"not null"`
	Code      string    `json:"code" gorm:"not null"`
	Status    OTPStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Session struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type RequestOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type AuthResponse struct {
	User  any    `json:"user"`
	Token string `json:"token"`
}

func (OTP) TableName() string {
	return "otps"
}

func (Session) TableName() string {
	return "sessions"
}
