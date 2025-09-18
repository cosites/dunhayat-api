package adapter

import (
	"context"

	"dunhayat-api/internal/auth/port"
	"dunhayat-api/internal/users/repository"

	"github.com/google/uuid"
)

type MiddlewareUserAdapter struct {
	userRepo repository.UserRepository
}

func NewMiddlewareUserAdapter(
	userRepo repository.UserRepository,
) port.UserReader {
	return &MiddlewareUserAdapter{
		userRepo: userRepo,
	}
}

func (s *MiddlewareUserAdapter) GetUserByID(
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
