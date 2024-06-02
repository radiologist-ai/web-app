package patientrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
	"github.com/rs/zerolog"
)

type PatientRepo struct {
	db     *sqlx.DB
	logger *zerolog.Logger
}

func (pr *PatientRepo) GetSelfPatientIDOfUser(ctx context.Context, userID int) (uuid.UUID, error) {
	var res uuid.UUID
	query := `SELECT id FROM patients WHERE user_id = $1 AND creator_id = $1`
	err := pr.db.QueryRowxContext(ctx, query, userID).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, fmt.Errorf("%w%w", customerrors.NotFoundError, err)
		}
		pr.logger.Error().Err(err).Int("userID", userID).Msg("error selecting patient")
		return uuid.UUID{}, fmt.Errorf("%w%w", customerrors.InternalErrorSQL, err)
	}
	return res, nil

}

func (pr *PatientRepo) LinkPatient(ctx context.Context, user domain.UserRepoModel, patientID uuid.UUID) error {
	query := `UPDATE patients SET user_id = $1 WHERE id = $2 AND user_id IS NULL`
	res, err := pr.db.ExecContext(ctx, query, user.ID, patientID)
	if err != nil {
		return fmt.Errorf("%w%w", customerrors.InternalErrorSQL, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w%w", customerrors.InternalError, err)
	}
	if affected == 0 {
		return customerrors.NotFoundError
	}
	return nil
}

func (pr *PatientRepo) SelectAllByUser(ctx context.Context, userID int) ([]domain.PatientAccountInfo, error) {
	res := make([]domain.PatientAccountInfo, 0)
	query := `
		SELECT 
		    p.id as id, 
		    p.name as name, 
		    p.patient_identifier as patient_identifier, 
		    u.first_name || ' ' || u.last_name as doctor_full_name
		FROM patients AS p INNER JOIN 
		    users AS u ON u.id = p.creator_id
		WHERE p.user_id = $1 AND p.creator_id != $1
`
	err := pr.db.SelectContext(ctx, &res, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%w%w", customerrors.InternalErrorSQL, err)
	}
	return res, nil
}

func (pr *PatientRepo) SelectOne(ctx context.Context, userID int, patientID uuid.UUID) (domain.PatientRepoModel, error) {
	var res domain.PatientRepoModel
	query := `SELECT id, user_id, creator_id, name, patient_identifier, created_at, updated_at
			  FROM patients WHERE creator_id = $1 AND id = $2`
	err := pr.db.QueryRowxContext(ctx, query, userID, patientID).StructScan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PatientRepoModel{}, fmt.Errorf("%w%w", customerrors.NotFoundError, err)
		}
		pr.logger.Error().Err(err).Any("patientID", patientID).Int("doctorID", userID).Msg("error selecting patient")
		return domain.PatientRepoModel{}, fmt.Errorf("%w%w", customerrors.InternalErrorSQL, err)
	}
	return res, nil
}

func (pr *PatientRepo) InsertPatient(ctx context.Context, patient domain.PatientRepoModel) (domain.PatientRepoModel, error) {
	var res domain.PatientRepoModel
	query := `
		INSERT INTO patients (creator_id, name, patient_identifier) 
		VALUES ($1, $2, $3)
		RETURNING id, user_id, creator_id, name, patient_identifier, created_at, updated_at
`
	if err := pr.db.QueryRowxContext(ctx, query, patient.CreatorID, patient.Name, patient.PatientIdentifier).StructScan(&res); err != nil {
		return domain.PatientRepoModel{}, fmt.Errorf("%w%w", customerrors.InternalErrorSQL, err)
	}
	return res, nil
}

func (pr *PatientRepo) SelectAll(ctx context.Context, userID int) ([]domain.PatientRepoModel, error) {
	res := make([]domain.PatientRepoModel, 0)
	query := `
		SELECT id, user_id, creator_id, name, patient_identifier, created_at, updated_at
		FROM patients WHERE creator_id=$1 
		ORDER BY updated_at DESC
`
	err := pr.db.SelectContext(ctx, &res, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%w%w", customerrors.InternalErrorSQL, err)
	}
	return res, nil
}

func New(db *sqlx.DB, logger *zerolog.Logger) (*PatientRepo, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}
	if db == nil {
		return nil, errors.New("db required")
	}
	return &PatientRepo{
		db:     db,
		logger: logger,
	}, nil
}
