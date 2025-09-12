package adapter

import (
	"context"

	"dunhayat-api/internal/auth/port"
	"dunhayat-api/internal/users/repository"

	"github.com/google/uuid"
)

type MiddlewareUserService struct {
	userRepo repository.UserRepository
}

func NewMiddlewareUserService(
	userRepo repository.UserRepository,
) port.UserReader {
	return &MiddlewareUserService{
		userRepo: userRepo,
	}
}

func (s *MiddlewareUserService) GetUserByID(
	ctx context.Context,
	userID uuid.UUID,
) (*port.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	return &port.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Email:     user.Email,
		Verified:  int(user.Verified),
		LastLogin: nil,
	}, nil
}
