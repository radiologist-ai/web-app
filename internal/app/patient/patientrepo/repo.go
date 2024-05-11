package patientrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
	"github.com/rs/zerolog"
)

type PatientRepo struct {
	db     *sqlx.DB
	logger *zerolog.Logger
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
