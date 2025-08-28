package redis

import (
	"context"
	"fmt"
	"time"

	"dunhayat-api/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func Connect(cfg *Config, log *logger.Logger) (*redis.Client, error) {
	log.Info(
		"Connecting to Redis",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.Int("db", cfg.DB),
	)

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Error(
			"Failed to connect to Redis",
			zap.Error(err),
		)
		return nil, fmt.Errorf(
			"failed to connect to Redis: %w",
			err,
		)
	}

	log.Info(
		"Redis connection established successfully",
		zap.Int("poolSize", 10),
		zap.Int("minIdleConns", 5),
		zap.Duration("dialTimeout", 10*time.Second),
		zap.Duration("readTimeout", 30*time.Second),
		zap.Duration("writeTimeout", 30*time.Second),
	)

	return client, nil
}

func Close(client *redis.Client) error {
	return client.Close()
}
