package usersrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
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
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("start transaction")
		return user, fmt.Errorf("%w%w", customerrors.InternalErrorSQL, err)
	}
	defer func() {
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				r.logger.Error().Err(errTx).Msg("rollback transaction")
			}
		}
	}()
	if user, err = r.insertUser(ctx, tx, model); err != nil {
		return
	}
	if !model.IsDoctor {
		if err = r.insertSelfPatient(ctx, tx, user.ID, fmt.Sprintf("%s %s", user.FirstName, user.LastName)); err != nil {
			return
		}
	}
	if err = tx.Commit(); err != nil {
		return domain.UserRepoModel{}, fmt.Errorf("%w%w", customerrors.InternalErrorSQL, err)
	}

	return user, nil
}

func (r *Repo) insertUser(ctx context.Context, tx *sqlx.Tx, model domain.UserRepoModel) (user domain.UserRepoModel, err error) {
	q := `Insert into users
    (email, first_name, last_name, password_hash, is_doctor) values 
    ($1, $2, $3, $4, $5)
    RETURNING id, email, first_name, last_name, password_hash, is_doctor, created_at, updated_at`
	if err = tx.QueryRowxContext(ctx, q, model.Email, model.FirstName, model.LastName, model.PasswordHash, model.IsDoctor).StructScan(&user); err != nil {
		r.logger.Error().Err(err).Msg("insert user")
		return domain.UserRepoModel{}, fmt.Errorf("%w%w", customerrors.InternalErrorSQL, err)
	}
	return user, nil
}

func (r *Repo) insertSelfPatient(ctx context.Context, tx *sqlx.Tx, userID int, name string) error {
	q := `INSERT INTO patients (user_id, creator_id, name) VALUES ($1, $2, $3)`
	_, err := tx.ExecContext(ctx, q, userID, userID, name)
	if err != nil {
		return fmt.Errorf("%w%w", customerrors.InternalErrorSQL, err)
	}
	return nil
}
