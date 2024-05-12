package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type (
	PatientRepoModel struct {
		ID                uuid.UUID `db:"id"`
		UserID            *int      `db:"user_id"`
		CreatorID         int       `db:"creator_id"`
		Name              string    `db:"name"`
		PatientIdentifier string    `db:"patient_identifier"`
		CreatedAt         time.Time `db:"created_at"`
		UpdatedAt         time.Time `db:"updated_at"`
	}

	PatientAccountInfo struct {
		ID                uuid.UUID `db:"id"`
		Name              string    `db:"name"`
		PatientIdentifier string    `db:"patient_identifier"`
		DoctorFullName    string    `db:"doctor_full_name"`
	}
)

type (
	PatientService interface {
		CreatePatient(ctx context.Context, creator UserRepoModel, form PatientRepoModel) (PatientRepoModel, error)
		GetAll(ctx context.Context, currentUser UserRepoModel) ([]PatientRepoModel, error)
		GetOne(ctx context.Context, currentUser UserRepoModel, patientID uuid.UUID) (PatientRepoModel, error)
		LinkPatient(ctx context.Context, user UserRepoModel, patientID string) error
		GetAllByUser(ctx context.Context, user UserRepoModel) ([]PatientAccountInfo, error)
	}

	PatientRepository interface {
		InsertPatient(ctx context.Context, patient PatientRepoModel) (PatientRepoModel, error)
		SelectAll(ctx context.Context, userID int) ([]PatientRepoModel, error)
		SelectOne(ctx context.Context, userID int, patientID uuid.UUID) (PatientRepoModel, error)
		SelectAllByUser(ctx context.Context, userID int) ([]PatientAccountInfo, error)
		LinkPatient(ctx context.Context, user UserRepoModel, patientID uuid.UUID) error
	}
)
