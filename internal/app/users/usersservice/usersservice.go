package usersservice

import (
	"context"
	"errors"
	"github.com/radiologist-ai/web-app/internal/domain"
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

func (s *Service) CreateOne(ctx context.Context, user domain.UserForm) (domain.UserRepoModel, error) {
	if err := s.validateRegisterForm(user); err != nil {
		return domain.UserRepoModel{}, err
	}
	repoModel := s.userFormToUserRepoModel(user)
	newUser, err := s.repo.InsertOne(ctx, repoModel)
	if err != nil {
		return domain.UserRepoModel{}, err
	}
	return newUser, nil
}
