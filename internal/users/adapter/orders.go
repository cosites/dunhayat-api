package adapter

import (
	"context"

	"dunhayat-api/internal/orders/port"
	"dunhayat-api/internal/users/repository"

	"github.com/google/uuid"
)

type OrdersUserAdapter struct {
	userRepo repository.UserRepository
}

func NewOrdersUserAdapter(
	userRepo repository.UserRepository,
) port.UserPort {
	return &OrdersUserAdapter{
		userRepo: userRepo,
	}
}

func (s *OrdersUserAdapter) GetUserByID(
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
	}, nil
}
