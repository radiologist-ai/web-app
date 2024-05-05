package usersservice

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
	"github.com/rs/zerolog"
	"time"
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

func (s *Service) GenerateToken(secret []byte, email string) (token string, err error) {
	payload := jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(time.Hour * 24 * 365).Unix(),
	}

	jwToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err = jwToken.SignedString(secret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) ValidateToken(secret []byte, token string) (email string, err error) {
	jwToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) { return secret, nil })
	if err != nil {
		return "", err
	}
	claims, ok := jwToken.Claims.(jwt.MapClaims)
	if !ok || !jwToken.Valid {
		return "", customerrors.ValidationErrorJWT
	}
	email, err = claims.GetSubject()
	if err != nil {
		return "", err
	}
	return email, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (user domain.UserRepoModel, ok bool, err error) {
	if user, ok, err = s.repo.SelectByEmail(ctx, email); err != nil {
		return domain.UserRepoModel{}, false, err
	}
	return
}

func (s *Service) CreateOne(ctx context.Context, user domain.UserForm) (domain.UserRepoModel, error) {
	if err := s.validateRegisterForm(user); err != nil {
		return domain.UserRepoModel{}, err
	}
	repoModel, err := s.userFormToUserRepoModel(user)
	if err != nil {
		return domain.UserRepoModel{}, err
	}
	newUser, err := s.repo.InsertOne(ctx, repoModel)
	if err != nil {
		return domain.UserRepoModel{}, err
	}
	return newUser, nil
}
