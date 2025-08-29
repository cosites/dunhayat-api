package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"dunhayat-api/internal/auth"
	"dunhayat-api/pkg/logger"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisOTPRepository struct {
	client *redis.Client
	logger logger.Interface
}

func NewRedisOTPRepository(
	client *redis.Client,
	logger logger.Interface,
) OTPRepository {
	return &RedisOTPRepository{
		client: client,
		logger: logger,
	}
}

func (r *RedisOTPRepository) Create(
	ctx context.Context,
	otp *auth.OTP,
) error {
	key := fmt.Sprintf("otp:%s", otp.Phone)

	r.logger.Debug(
		"Creating OTP in Redis",
		zap.String("phone", otp.Phone),
		zap.String("key", key),
		zap.Time("expiresAt", otp.ExpiresAt),
	)

	otpData, err := json.Marshal(otp)
	if err != nil {
		r.logger.Error(
			"Failed to marshal OTP",
			zap.Error(err),
			zap.String("phone", otp.Phone),
		)
		return fmt.Errorf(
			"failed to marshal OTP: %w",
			err,
		)
	}

	expiration := time.Until(otp.ExpiresAt)
	if expiration <= 0 {
		r.logger.Warn(
			"OTP already expired",
			zap.String("phone", otp.Phone),
			zap.Time("expiresAt", otp.ExpiresAt),
		)
		return fmt.Errorf("OTP already expired")
	}

	err = r.client.Set(ctx, key, otpData, expiration).Err()
	if err != nil {
		r.logger.Error(
			"Failed to store OTP in Redis",
			zap.Error(err),
			zap.String("phone", otp.Phone),
			zap.String("key", key),
		)
		return fmt.Errorf(
			"failed to store OTP in Redis: %w",
			err,
		)
	}

	r.logger.Debug(
		"OTP stored successfully in Redis",
		zap.String("phone", otp.Phone),
		zap.Duration("expiration", expiration),
	)
	return nil
}

func (r *RedisOTPRepository) GetByPhone(
	ctx context.Context,
	phone string,
) (*auth.OTP, error) {
	key := fmt.Sprintf("otp:%s", phone)

	otpData, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf(
			"failed to get OTP from Redis: %w",
			err,
		)
	}

	var otp auth.OTP
	if err := json.Unmarshal([]byte(otpData), &otp); err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal OTP: %w",
			err,
		)
	}

	return &otp, nil
}

func (r *RedisOTPRepository) Update(
	ctx context.Context,
	otp *auth.OTP,
) error {
	key := fmt.Sprintf("otp:%s", otp.Phone)

	otpData, err := json.Marshal(otp)
	if err != nil {
		return fmt.Errorf(
			"failed to marshal OTP: %w",
			err,
		)
	}

	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to get TTL: %w", err)
	}

	if ttl > 0 {
		err = r.client.Set(ctx, key, otpData, ttl).Err()
		if err != nil {
			return fmt.Errorf(
				"failed to update OTP in Redis: %w",
				err,
			)
		}
	} else {
		expiration := time.Until(otp.ExpiresAt)
		if expiration > 0 {
			err = r.client.Set(ctx, key, otpData, expiration).Err()
			if err != nil {
				return fmt.Errorf(
					"failed to create OTP in Redis: %w",
					err,
				)
			}
		}
	}

	return nil
}

func (r *RedisOTPRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r *RedisOTPRepository) CleanExpired(ctx context.Context) error {
	return nil
}

func (r *RedisOTPRepository) GetTTL(
	ctx context.Context,
	phone string,
) (time.Duration, error) {
	key := fmt.Sprintf("otp:%s", phone)

	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL: %w", err)
	}

	return ttl, nil
}

func (r *RedisOTPRepository) InvalidateOTP(
	ctx context.Context,
	phone string,
) error {
	key := fmt.Sprintf("otp:%s", phone)

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to invalidate OTP: %w", err)
	}

	return nil
}

func (r *RedisOTPRepository) GetOTPCount(
	ctx context.Context,
	phone string,
) (int64, error) {
	pattern := fmt.Sprintf("otp:%s*", phone)

	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get OTP keys: %w", err)
	}

	return int64(len(keys)), nil
}
