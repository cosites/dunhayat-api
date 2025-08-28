package db

import (
	"time"

	"dunhayat-api/pkg/config"
	"dunhayat-api/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func Connect(
	dbConfig *config.DatabaseConfig,
	log *logger.Logger,
) (*gorm.DB, error) {
	log.Info(
		"Connecting to database",
		zap.String("host", dbConfig.Host),
		zap.Int("port", dbConfig.Port),
		zap.String("database", dbConfig.DBName),
	)

	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	}

	db, err := gorm.Open(
		postgres.Open(
			dbConfig.GetDSN(),
		),
		gormConfig,
	)
	if err != nil {
		log.Error(
			"Failed to connect to database",
			zap.Error(err),
		)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error(
			"Failed to get underlying sql.DB",
			zap.Error(err),
		)
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info(
		"Database connection established successfully",
		zap.Int("maxIdleConns", 10),
		zap.Int("maxOpenConns", 100),
		zap.Duration("connMaxLifetime", time.Hour),
	)

	return db, nil
}
