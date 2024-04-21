package usersservice

import (
	"ai-radiologist/internal/domain"
	"context"
	"errors"
	"github.com/rs/zerolog"
)

type Service struct {
	logger *zerolog.Logger
	repo   domain.UsersRepository
}

func New(logger *zerolog.Logger, repo domain.UsersRepository) (*Service, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}
	if repo == nil {
		return nil, errors.New("repo required")
	}
	return &Service{logger: logger, repo: repo}, nil
}

func (s *Service) GenerateToken(ctx context.Context, email string) (token string, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) ValidateToken(ctx context.Context, token string) (email string, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetByEmail(ctx context.Context, email string) (user domain.UserRepoModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) CreateOne(ctx context.Context, user domain.UserRepoModel) (domain.UserRepoModel, error) {
	//TODO implement me
	panic("implement me")
}
