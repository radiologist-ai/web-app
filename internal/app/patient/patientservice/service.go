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
	s3     domain.S3Adapter
}

func (ps *PatientService) GetSelfPatientIDOfUser(ctx context.Context, user domain.UserRepoModel) (uuid.UUID, error) {
	if user.IsDoctor {
		return uuid.UUID{}, customerrors.NeedToBePatient
	}
	patientID, err := ps.repo.GetSelfPatientIDOfUser(ctx, user.ID)
	if err != nil {
		ps.logger.Error().Err(err).Any("user", user).Ctx(ctx).Msg("patientService.GetSelfPatientIDOfUser")
		return uuid.UUID{}, err
	}
	return patientID, nil
}

func (ps *PatientService) ListAIGeneratedReportsByPatientID(ctx context.Context, patientID uuid.UUID) ([]domain.ReportForPatient, error) {
	result, err := ps.repo.SelectAIGeneratedReportsByPatient(ctx, patientID)
	if err != nil {
		ps.logger.Error().Err(err).Any("patientID", patientID).Ctx(ctx).Msg("patientService.ListAIGeneratedReportsByPatientID")
		return nil, err
	}
	for i := range result {
		result[i].DoctorName = "Radiologist AI"
		result[i].ImagePath = ps.s3.GetPublicLink(result[i].ImagePath)
	}
	return result, nil
}

func (ps *PatientService) GetApprovedReportsByUser(ctx context.Context, user domain.UserRepoModel) ([]domain.ReportForPatient, error) {
	if user.IsDoctor {
		return make([]domain.ReportForPatient, 0), customerrors.NeedToBePatient
	}
	result, err := ps.repo.SelectApprovedReportsByUser(ctx, user.ID)
	if err != nil {
		ps.logger.Error().Err(err).Any("user", user).Ctx(ctx).Msg("patientService.GetApprovedReportsByUser")
		return nil, err
	}
	for i := range result {
		result[i].ImagePath = ps.s3.GetPublicLink(result[i].ImagePath)
	}
	return result, nil
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

func New(logger *zerolog.Logger, repo domain.PatientRepository, s3 domain.S3Adapter) (*PatientService, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger can not be nil")
	}
	if repo == nil {
		return nil, fmt.Errorf("repo can not be nil")
	}
	if s3 == nil {
		return nil, fmt.Errorf("s3 can not be nil")
	}
	return &PatientService{
		repo:   repo,
		logger: logger,
		s3:     s3,
	}, nil
}
