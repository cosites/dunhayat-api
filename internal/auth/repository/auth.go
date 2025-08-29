package repository

import (
	"context"
	"errors"
	"time"

	"dunhayat-api/internal/auth"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OTPRepository interface {
	Create(ctx context.Context, otp *auth.OTP) error
	GetByPhone(ctx context.Context, phone string) (*auth.OTP, error)
	Update(ctx context.Context, otp *auth.OTP) error
	Delete(ctx context.Context, id uuid.UUID) error
	CleanExpired(ctx context.Context) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *auth.Session) error
	GetByToken(ctx context.Context, token string) (*auth.Session, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CleanExpired(ctx context.Context) error
}

type postgresOTPRepository struct {
	db *gorm.DB
}

type postgresSessionRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) OTPRepository {
	return &postgresOTPRepository{db: db}
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &postgresSessionRepository{db: db}
}

func (r *postgresOTPRepository) Create(
	ctx context.Context,
	otp *auth.OTP,
) error {
	return r.db.WithContext(ctx).Create(otp).Error
}

func (r *postgresOTPRepository) GetByPhone(
	ctx context.Context,
	phone string,
) (*auth.OTP, error) {
	var otp auth.OTP
	err := r.db.WithContext(ctx).Where(
		"phone = ?",
		phone,
	).First(&otp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &otp, nil
}

func (r *postgresOTPRepository) Update(
	ctx context.Context,
	otp *auth.OTP,
) error {
	return r.db.WithContext(ctx).Save(otp).Error
}

func (r *postgresOTPRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return r.db.WithContext(ctx).Delete(&auth.OTP{}, id).Error
}

func (r *postgresOTPRepository) CleanExpired(
	ctx context.Context,
) error {
	return r.db.WithContext(ctx).Where(
		"expires_at < ?",
		time.Now(),
	).Delete(&auth.OTP{}).Error
}

func (r *postgresSessionRepository) Create(
	ctx context.Context,
	session *auth.Session,
) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *postgresSessionRepository) GetByToken(
	ctx context.Context,
	token string,
) (*auth.Session, error) {
	var session auth.Session
	err := r.db.WithContext(ctx).Where(
		"token = ?",
		token,
	).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *postgresSessionRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return r.db.WithContext(ctx).Delete(&auth.Session{}, id).Error
}

func (r *postgresSessionRepository) CleanExpired(
	ctx context.Context,
) error {
	return r.db.WithContext(ctx).Where(
		"expires_at < ?",
		time.Now(),
	).Delete(&auth.Session{}).Error
}
