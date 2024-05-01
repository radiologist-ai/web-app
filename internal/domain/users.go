package domain

import (
	"context"
	"time"
)

type UsersService interface {
	AuthService
	UsersValidator
	GetByEmail(ctx context.Context, email string) (user UserRepoModel, err error)
	CreateOne(ctx context.Context, user UserForm) (UserRepoModel, error)
}

type UsersRepository interface {
	SelectByEmail(ctx context.Context, email string) (user UserRepoModel, err error)
	InsertOne(ctx context.Context, model UserRepoModel) (user UserRepoModel, err error)
}

type AuthService interface {
	GenerateToken(ctx context.Context, email string) (token string, err error)
	ValidateToken(ctx context.Context, token string) (email string, err error)
}

type UsersValidator interface {
	ValidatePassword(password string) error
}

type UserRepoModel struct {
	ID           int       `db:"id"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password_hash"`
	IsDoctor     bool      `db:"is_doctor"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type UserForm struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	IsDoctor  bool
}
