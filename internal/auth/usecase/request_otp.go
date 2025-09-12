package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"time"
	"unicode"

	"dunhayat-api/internal/auth"
	"dunhayat-api/internal/auth/repository"
	"dunhayat-api/pkg/logger"
	"dunhayat-api/pkg/sms"

	"go.uber.org/zap"
)

type RequestOTPUseCase interface {
	Execute(ctx context.Context, phone string) (*auth.OTP, error)
}

type requestOTPUseCase struct {
	otpRepo     repository.OTPRepository
	smsProvider sms.Provider
	template    string
	logger      logger.Interface
}

func NewRequestOTPUseCase(
	otpRepo repository.OTPRepository,
	smsProvider sms.Provider,
	template string,
	logger logger.Interface,
) RequestOTPUseCase {
	return &requestOTPUseCase{
		otpRepo:     otpRepo,
		smsProvider: smsProvider,
		template:    template,
		logger:      logger,
	}
}

func (uc *requestOTPUseCase) Execute(
	ctx context.Context,
	phone string,
) (*auth.OTP, error) {
	uc.logger.Info("Starting OTP request", zap.String("phone", phone))

	if err := uc.validatePhone(phone); err != nil {
		uc.logger.Error("Invalid phone number format", zap.Error(err))
		return nil, fmt.Errorf("invalid phone number format: %w", err)
	}

	if err := uc.checkRateLimit(ctx, phone); err != nil {
		uc.logger.Warn(
			"Rate limit exceeded",
			zap.String("phone", phone),
			zap.Error(err),
		)
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	otpCode := uc.generateOTP()
	uc.logger.Info("Generated OTP code", zap.String("otp_code", otpCode))

	expiresAt := time.Now().Add(10 * time.Minute)
	uc.logger.Info("OTP expires at", zap.Time("expires_at", expiresAt))

	otp := &auth.OTP{
		Phone:     phone,
		Code:      otpCode,
		Status:    auth.OTPStatusPending,
		ExpiresAt: expiresAt,
	}

	uc.logger.Info("Saving OTP to repository...")
	if err := uc.otpRepo.Create(ctx, otp); err != nil {
		uc.logger.Error("Failed to save OTP", zap.Error(err))
		return nil, fmt.Errorf("failed to save OTP: %w", err)
	}
	uc.logger.Info("OTP saved successfully to repository")

	uc.logger.Info(
		"Sending OTP via SMS using template",
		zap.String("template", uc.template),
	)

	if err := uc.smsProvider.SendOTP(
		ctx, phone, otpCode, uc.template,
	); err != nil {
		uc.logger.Error("Failed to send OTP SMS", zap.Error(err))
		otp.Status = auth.OTPStatusFailed
		_ = uc.otpRepo.Update(ctx, otp)
		return nil, fmt.Errorf("failed to send OTP SMS: %w", err)
	}

	uc.logger.Info("OTP sent successfully via SMS")
	return otp, nil
}

func (uc *requestOTPUseCase) generateOTP() string {
	digits := make([]byte, 6)
	randomBytes := make([]byte, 6)

	if _, err := rand.Read(randomBytes); err != nil {
		for i := range digits {
			digits[i] = byte('0' + (time.Now().UnixNano() % 10))
		}
		return string(digits)
	}

	for i := range digits {
		digits[i] = byte('0' + (randomBytes[i] % 10))
	}

	return string(digits)
}

func (uc *requestOTPUseCase) validatePhone(phone string) error {
	if len(phone) == 11 {
		if !strings.HasPrefix(phone, "09") {
			return fmt.Errorf(
				"phone number shall start with either 09 or +98",
			)
		}
		phone = "+98" + phone[1:]
	}

	if len(phone) != 13 && !strings.HasPrefix(phone, "+98") {
		return fmt.Errorf(
			"phone number shall start with either 09 or +98",
		)
	}

	for _, char := range phone[3:] {
		if !unicode.IsNumber(char) {
			return fmt.Errorf(
				"phone number shall contain only digits after country code",
			)
		}
	}

	return nil
}

func (uc *requestOTPUseCase) checkRateLimit(
	ctx context.Context,
	phone string,
) error {
	existingOTP, err := uc.otpRepo.GetByPhone(ctx, phone)
	if err != nil {
		uc.logger.Error("Failed to check existing OTP", zap.Error(err))
		return fmt.Errorf("failed to check rate limit: %w", err)
	}

	if existingOTP != nil {
		timeSinceLastOTP := time.Since(existingOTP.CreatedAt)
		if timeSinceLastOTP < 2*time.Minute {
			remainingTime := 2*time.Minute - timeSinceLastOTP
			return fmt.Errorf(
				"please wait %v before requesting another OTP",
				remainingTime.Round(time.Second),
			)
		}
	}

	return nil
}
