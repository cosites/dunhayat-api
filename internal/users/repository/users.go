package repository

import (
	"context"
	"errors"

	"dunhayat-api/internal/users"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *users.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*users.User, error)
	GetByPhone(ctx context.Context, phone string) (*users.User, error)
	Update(ctx context.Context, user *users.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type AddressRepository interface {
	Create(ctx context.Context, address *users.Address) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]users.Address, error)
	Update(ctx context.Context, address *users.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type postgresUserRepository struct {
	db *gorm.DB
}

type postgresAddressRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &postgresAddressRepository{db: db}
}

func (r *postgresUserRepository) Create(
	ctx context.Context,
	user *users.User,
) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *postgresUserRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*users.User, error) {
	var user users.User
	err := r.db.WithContext(ctx).Where(
		"id = ?", id,
	).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *postgresUserRepository) GetByPhone(
	ctx context.Context,
	phone string,
) (*users.User, error) {
	var user users.User
	err := r.db.WithContext(ctx).Where(
		"phone = ?", phone,
	).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *postgresUserRepository) Update(
	ctx context.Context,
	user *users.User,
) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *postgresUserRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return r.db.WithContext(ctx).Delete(&users.User{}, id).Error
}

func (r *postgresAddressRepository) Create(
	ctx context.Context,
	address *users.Address,
) error {
	return r.db.WithContext(ctx).Create(address).Error
}

func (r *postgresAddressRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]users.Address, error) {
	var addresses []users.Address
	err := r.db.WithContext(ctx).Where(
		"user_id = ?",
		userID,
	).Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *postgresAddressRepository) Update(
	ctx context.Context,
	address *users.Address,
) error {
	return r.db.WithContext(ctx).Save(address).Error
}

func (r *postgresAddressRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return r.db.WithContext(ctx).Delete(&users.Address{}, id).Error
}
