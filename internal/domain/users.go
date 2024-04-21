package domain

import "context"

type UsersService interface {
	AuthService
	GetByEmail(ctx context.Context, email string) (user UserRepoModel, err error)
	CreateOne(ctx context.Context, user UserRepoModel) (UserRepoModel, error)
}

type UsersRepository interface {
	SelectByEmail(ctx context.Context, email string) (user UserRepoModel, err error)
	InsertOne(ctx context.Context, model UserRepoModel) (user UserRepoModel, err error)
}

type AuthService interface {
	GenerateToken(ctx context.Context, email string) (token string, err error)
	ValidateToken(ctx context.Context, token string) (email string, err error)
}

type UserRepoModel struct {
	ID           int
	FirstName    string
	LastName     string
	Email        string
	PasswordHash []byte
}

type UserForm struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}
