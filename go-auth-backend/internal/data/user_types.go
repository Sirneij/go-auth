package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"goauthbackend.johnowolabiidogun.dev/internal/types"
)

type UserProfile struct {
	ID          *uuid.UUID     `json:"id"`
	UserID      *uuid.UUID     `json:"user_id"`
	PhoneNumber *string        `json:"phone_number"`
	BirthDate   types.NullTime `json:"birth_date"`
	GithubLink  *string        `json:"github_link"`
}

type User struct {
	ID          uuid.UUID   `json:"id"`
	Email       string      `json:"email"`
	Password    password    `json:"-"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	IsActive    bool        `json:"is_active"`
	IsStaff     bool        `json:"is_staff"`
	IsSuperuser bool        `json:"is_superuser"`
	Thumbnail   *string     `json:"thumbnail"`
	DateJoined  time.Time   `json:"date_joined"`
	Profile     UserProfile `json:"profile"`
}

type password struct {
	plaintext *string
	hash      string
}

type UserModel struct {
	DB *sql.DB
}

type UserID struct {
	Id uuid.UUID
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)
