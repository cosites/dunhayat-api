package adapter

import (
	"context"
	"time"

	"dunhayat-api/internal/auth/port"
	"dunhayat-api/internal/users"
	"dunhayat-api/internal/users/repository"

	"github.com/google/uuid"
)

type AuthUserAdapter struct {
	userRepo    repository.UserRepository
	addressRepo repository.AddressRepository
}

func NewAuthUserAdapter(
	userRepo repository.UserRepository,
	addressRepo repository.AddressRepository,
) port.UserPort {
	return &AuthUserAdapter{
		userRepo:    userRepo,
		addressRepo: addressRepo,
	}
}

func (s *AuthUserAdapter) FindUserByPhone(
	ctx context.Context,
	phone string,
) (*port.User, error) {
	user, err := s.userRepo.GetByPhone(ctx, phone)
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

func (s *AuthUserAdapter) CreateUser(
	ctx context.Context,
	user *port.User,
) error {
	domainUser := &users.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Email:     user.Email,
		Verified:  users.VerificationLevel(user.Verified),
		LastLogin: nil,
	}

	return s.userRepo.Create(ctx, domainUser)
}

func (s *AuthUserAdapter) UpdateUserLastLogin(
	ctx context.Context,
	userID uuid.UUID,
) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}

	now := time.Now()
	user.LastLogin = &now

	return s.userRepo.Update(ctx, user)
}

func (s *AuthUserAdapter) CreateAddress(
	ctx context.Context,
	address *port.Address,
) error {
	domainAddress := &users.Address{
		ID:         address.ID,
		UserID:     address.UserID,
		Address:    address.Address,
		PostalCode: address.PostalCode,
	}

	return s.addressRepo.Create(ctx, domainAddress)
}

func (s *AuthUserAdapter) GetUserAddresses(
	ctx context.Context,
	userID uuid.UUID,
) ([]port.Address, error) {
	domainAddresses, err := s.addressRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	addresses := make([]port.Address, len(domainAddresses))
	for i, addr := range domainAddresses {
		addresses[i] = port.Address{
			ID:         addr.ID,
			UserID:     addr.UserID,
			Address:    addr.Address,
			PostalCode: addr.PostalCode,
		}
	}

	return addresses, nil
}
