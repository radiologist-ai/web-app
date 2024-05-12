package patientservice

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/radiologist-ai/web-app/internal/domain"
	"github.com/radiologist-ai/web-app/internal/domain/customerrors"
	"github.com/rs/zerolog"
)

type PatientService struct {
	repo   domain.PatientRepository
	logger *zerolog.Logger
}

func (ps *PatientService) GetAllByUser(ctx context.Context, user domain.UserRepoModel) ([]domain.PatientAccountInfo, error) {
	if user.IsDoctor {
		return make([]domain.PatientAccountInfo, 0), customerrors.NeedToBePatient
	}
	result, err := ps.repo.SelectAllByUser(ctx, user.ID)
	if err != nil {
		ps.logger.Error().Err(err).Any("user", user).Ctx(ctx).Msg("patientService.GetAllByUser")
		return nil, err
	}
	return result, nil
}

func (ps *PatientService) LinkPatient(ctx context.Context, user domain.UserRepoModel, patientID string) error {
	patientUUID, err := uuid.Parse(patientID)
	if err != nil {
		return customerrors.ValidationErrorUUID
	}
	if err := ps.repo.LinkPatient(ctx, user, patientUUID); err != nil {
		ps.logger.Error().Err(err).Any("user", user).Str("patientID", patientID).Msg("failed to link patient")
		return err
	}
	return nil
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

func (ps *PatientService) GetAll(ctx context.Context, currentUser domain.UserRepoModel) ([]domain.PatientRepoModel, error) {
	if !currentUser.IsDoctor {
		return nil, fmt.Errorf("%wonly doctor can create patient. ", customerrors.NeedToBeDoctor)
	}
	res, err := ps.repo.SelectAll(ctx, currentUser.ID)
	if err != nil {
		ps.logger.Error().Err(err).Int("user_id", currentUser.ID).Msg("patient query failed")
		return nil, err
	}
	return res, nil
}

func (ps *PatientService) GetOne(ctx context.Context, currentUser domain.UserRepoModel, patientID uuid.UUID) (domain.PatientRepoModel, error) {
	if !currentUser.IsDoctor {
		return domain.PatientRepoModel{}, fmt.Errorf("%wonly doctor can create patient. ", customerrors.NeedToBeDoctor)
	}
	patient, err := ps.repo.SelectOne(ctx, currentUser.ID, patientID)
	if err != nil {
		ps.logger.Error().Err(err).Msg("patient query failed")
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
