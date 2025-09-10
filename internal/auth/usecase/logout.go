package usecase

import (
	"context"
	"fmt"

	authRepo "dunhayat-api/internal/auth/repository"
	"dunhayat-api/pkg/logger"

	"go.uber.org/zap"
)

type LogoutUseCase interface {
	Execute(ctx context.Context, token string) error
}

type logoutUseCase struct {
	sessionRepo authRepo.SessionRepository
	logger      logger.Interface
}

func NewLogoutUseCase(
	sessionRepo authRepo.SessionRepository,
	logger logger.Interface,
) LogoutUseCase {
	return &logoutUseCase{
		sessionRepo: sessionRepo,
		logger:      logger,
	}
}

func (uc *logoutUseCase) Execute(ctx context.Context, token string) error {
	uc.logger.Info(
		"Starting logout process",
		zap.String("token", token[:8]+"..."),
	)

	session, err := uc.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		uc.logger.Error("Failed to get session", zap.Error(err))
		return fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		uc.logger.Warn(
			"Session not found for logout",
			zap.String("token", token[:8]+"..."),
		)
		return nil
	}

	uc.logger.Info("Session found, proceeding with logout",
		zap.String("session_id", session.ID.String()),
		zap.String("user_id", session.UserID.String()))

	if err := uc.sessionRepo.DeleteByToken(ctx, token); err != nil {
		uc.logger.Error("Failed to delete session", zap.Error(err))
		return fmt.Errorf("failed to delete session: %w", err)
	}

	uc.logger.Info("Logout completed successfully",
		zap.String("session_id", session.ID.String()),
		zap.String("user_id", session.UserID.String()))

	return nil
}
