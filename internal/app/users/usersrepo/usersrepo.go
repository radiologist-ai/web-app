package usersrepo

import (
	"ai-radiologist/internal/domain"
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Repo struct {
	logger *zerolog.Logger
	db     *sqlx.DB
}

func New(logger *zerolog.Logger, db *sqlx.DB) (*Repo, error) {
	if logger == nil {
		return nil, errors.New("param logger is nil")
	}
	if db == nil {
		return nil, errors.New("param db is nil")
	}
	return &Repo{
		logger: logger,
		db:     db,
	}, nil
}

func (r *Repo) SelectByEmail(ctx context.Context, email string) (user domain.UserRepoModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) InsertOne(ctx context.Context, model domain.UserRepoModel) (user domain.UserRepoModel, err error) {
	//TODO implement me
	panic("implement me")
}
