package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"dunhayat-api/internal/auth"
	"dunhayat-api/internal/auth/port"
	authRepo "dunhayat-api/internal/auth/repository"
	"dunhayat-api/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type VerifyOTPUseCase interface {
	Execute(ctx context.Context, phone, code string) (*auth.AuthResponse, error)
}

type verifyOTPUseCase struct {
	otpRepo        authRepo.OTPRepository
	sessionRepo    authRepo.SessionRepository
	userService    port.UserService
	sessionTimeout time.Duration
	logger         logger.Interface
}

func NewVerifyOTPUseCase(
	otpRepo authRepo.OTPRepository,
	sessionRepo authRepo.SessionRepository,
	userService port.UserService,
	sessionTimeout time.Duration,
	logger logger.Interface,
) VerifyOTPUseCase {
	return &verifyOTPUseCase{
		otpRepo:        otpRepo,
		sessionRepo:    sessionRepo,
		userService:    userService,
		sessionTimeout: sessionTimeout,
		logger:         logger,
	}
}

func (uc *verifyOTPUseCase) Execute(
	ctx context.Context, phone, code string,
) (*auth.AuthResponse, error) {
	uc.logger.Info("Starting OTP verification", zap.String("phone", phone))

	otp, err := uc.getLatestValidOTP(ctx, phone)
	if err != nil {
		uc.logger.Error("Failed to get OTP", zap.Error(err))
		return nil, fmt.Errorf("failed to get OTP: %w", err)
	}

	if time.Now().After(otp.ExpiresAt) {
		uc.logger.Warn("OTP has expired", zap.String("phone", phone))
		otp.Status = auth.OTPStatusExpired
		uc.otpRepo.Update(ctx, otp)
		return nil, fmt.Errorf("OTP has expired")
	}

	if otp.Code != code {
		uc.logger.Warn("Invalid OTP code", zap.String("phone", phone))
		otp.Status = auth.OTPStatusFailed
		uc.otpRepo.Update(ctx, otp)
		return nil, fmt.Errorf("invalid OTP code")
	}

	uc.logger.Info("OTP validation successful", zap.String("phone", phone))

	otp.Status = auth.OTPStatusVerified
	if err := uc.otpRepo.Update(ctx, otp); err != nil {
		return nil, fmt.Errorf("failed to update OTP status: %w", err)
	}

	user, err := uc.getOrCreateUser(ctx, phone)
	if err != nil {
		uc.logger.Error("Failed to get or create user", zap.Error(err))
		return nil, fmt.Errorf("failed to get or create user: %w", err)
	}

	session, err := uc.createSession(ctx, user.ID)
	if err != nil {
		uc.logger.Error("Failed to create session", zap.Error(err))
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	if err := uc.userService.UpdateUserLastLogin(
		ctx,
		user.ID,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to update user last login: %w", err,
		)
	}

	addresses, err := uc.userService.GetUserAddresses(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get user addresses: %w", err,
		)
	}

	userData := map[string]any{
		"id":         user.ID,
		"phone":      user.Phone,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"verified":   user.Verified,
		"addresses":  addresses,
	}

	uc.logger.Info(
		"OTP verification completed successfully",
		zap.String("phone", phone),
		zap.String("user_id", user.ID.String()),
	)

	return &auth.AuthResponse{
		User:  userData,
		Token: session.Token,
	}, nil
}

func (uc *verifyOTPUseCase) getLatestValidOTP(
	ctx context.Context, phone string,
) (*auth.OTP, error) {
	otp, err := uc.otpRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to get OTP: %w", err)
	}

	if otp == nil {
		return nil, fmt.Errorf("no OTP found for phone number")
	}

	return otp, nil
}

func (uc *verifyOTPUseCase) getOrCreateUser(
	ctx context.Context, phone string,
) (*port.User, error) {
	user, err := uc.userService.FindUserByPhone(ctx, phone)
	if err == nil && user != nil {
		return user, nil
	}

	user = &port.User{
		ID:       uuid.New(),
		Phone:    phone,
		Verified: 0,
	}

	if err := uc.userService.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (uc *verifyOTPUseCase) createSession(
	ctx context.Context, userID uuid.UUID,
) (*auth.Session, error) {
	token := generateSessionToken()

	expiresAt := time.Now().Add(uc.sessionTimeout)

	session := &auth.Session{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func generateSessionToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)

	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		for i := range b {
			b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		}
		return string(b)
	}

	for i := range b {
		b[i] = charset[randomBytes[i]%byte(len(charset))]
	}

	return string(b)
}
