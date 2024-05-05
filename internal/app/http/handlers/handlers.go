package handlers

import (
	"errors"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/rs/zerolog"
)

type Handlers struct {
	logger *zerolog.Logger
	users  domain.UsersService
	secret []byte
}

func NewHandlers(logger *zerolog.Logger, users domain.UsersService, secret string) (*Handlers, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}
	if users == nil {
		return nil, errors.New("users is required")
	}
	if secret == "" {
		return nil, errors.New("secret is required")
	}
	return &Handlers{logger: logger, users: users, secret: []byte(secret)}, nil
}
