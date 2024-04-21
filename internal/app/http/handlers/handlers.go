package handlers

import (
	"ai-radiologist/internal/domain"
	"errors"
	"github.com/rs/zerolog"
)

type Handlers struct {
	logger *zerolog.Logger
	users  domain.UsersService
}

func NewHandlers(logger *zerolog.Logger, users domain.UsersService) (*Handlers, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}
	if users == nil {
		return nil, errors.New("users is required")
	}
	return &Handlers{logger: logger, users: users}, nil
}
