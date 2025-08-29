package port

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName *string   `json:"first_name,omitempty"`
	LastName  *string   `json:"last_name,omitempty"`
	Phone     string    `json:"phone"`
	Email     *string   `json:"email,omitempty"`
	Verified  int       `json:"verified"`
	LastLogin *string   `json:"last_login,omitempty"`
}

type Address struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Address    string    `json:"address"`
	PostalCode string    `json:"postal_code"`
}

type UserService interface {
	FindUserByPhone(ctx context.Context, phone string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUserLastLogin(ctx context.Context, userID uuid.UUID) error
	CreateAddress(ctx context.Context, address *Address) error
	GetUserAddresses(ctx context.Context, userID uuid.UUID) ([]Address, error)
}
