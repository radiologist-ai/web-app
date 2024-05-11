package patientservice

import (
	"context"
	"fmt"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
	"github.com/rs/zerolog"
)

type PatientService struct {
	repo   domain.PatientRepository
	logger *zerolog.Logger
}

func (ps *PatientService) CreatePatient(ctx context.Context, creator domain.UserRepoModel, form domain.PatientRepoModel) (domain.PatientRepoModel, error) {
	if !creator.IsDoctor {
		return domain.PatientRepoModel{}, fmt.Errorf("%wonly doctor can create patient. ", customerrors.NeedToBeDoctor)
	}
	form.CreatorID = creator.ID
	patient, err := ps.repo.InsertPatient(ctx, form)
	if err != nil {
		ps.logger.Error().Err(err).Msg("patient creation failed")
		return domain.PatientRepoModel{}, err
	}
	return patient, nil
}

func New(logger *zerolog.Logger, repo domain.PatientRepository) (*PatientService, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger can not be nil")
	}
	if repo == nil {
		return nil, fmt.Errorf("repo can not be nil")
	}
	return &PatientService{
		repo:   repo,
		logger: logger,
	}, nil
}
