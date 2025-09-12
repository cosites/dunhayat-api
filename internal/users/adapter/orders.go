package adapter

import (
	"context"

	"dunhayat-api/internal/orders/port"
	"dunhayat-api/internal/users/repository"

	"github.com/google/uuid"
)

type OrdersUserService struct {
	userRepo repository.UserRepository
}

func NewOrdersUserService(
	userRepo repository.UserRepository,
) port.UserService {
	return &OrdersUserService{
		userRepo: userRepo,
	}
}

func (s *OrdersUserService) GetUserByID(
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
