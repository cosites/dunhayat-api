package port

import (
	"context"

	"github.com/google/uuid"
)

type UserReader interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
}
