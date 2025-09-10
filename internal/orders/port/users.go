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
}

type UserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
}
