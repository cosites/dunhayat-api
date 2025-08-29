package users

import (
	"time"

	"github.com/google/uuid"
)

type VerificationLevel int

const (
	VerificationNone VerificationLevel = iota
	VerificationPhone
	VerificationEmail
	VerificationBoth
)

func (v VerificationLevel) String() string {
	switch v {
	case VerificationNone:
		return "none"
	case VerificationPhone:
		return "phone"
	case VerificationEmail:
		return "email"
	case VerificationBoth:
		return "both"
	default:
		return "unknown"
	}
}

type User struct {
	ID        uuid.UUID         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FirstName *string           `json:"first_name,omitempty"`
	LastName  *string           `json:"last_name,omitempty"`
	Phone     string            `json:"phone" gorm:"uniqueIndex;not null"`
	Email     *string           `json:"email,omitempty"`
	Verified  VerificationLevel `json:"verified" gorm:"type:smallint;default:0"`
	LastLogin *time.Time        `json:"last_login,omitempty"`
	CreatedAt time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
}

type Address struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Address    string    `json:"address" gorm:"type:text;not null"`
	PostalCode string    `json:"postal_code" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

func (Address) TableName() string {
	return "addresses"
}
