package sms

import "context"

type Provider interface {
	SendOTP(ctx context.Context, phone, code, template string) error
}
