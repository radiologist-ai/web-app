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

	ReportForPatient struct {
		ID          int       `db:"id"`
		PatientName string    `db:"name"`
		DoctorName  string    `db:"doctor_full_name"`
		ImagePath   string    `db:"image_path"`
		Text        string    `db:"report_text"`
		Approved    bool      `db:"approved"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
	}
)

type (
	PatientService interface {
		ListAIGeneratedReportsByPatientID(ctx context.Context, patientID uuid.UUID) ([]ReportForPatient, error)
		GetSelfPatientIDOfUser(ctx context.Context, user UserRepoModel) (uuid.UUID, error)
		GetApprovedReportsByUser(ctx context.Context, user UserRepoModel) ([]ReportForPatient, error)
		CreatePatient(ctx context.Context, creator UserRepoModel, form PatientRepoModel) (PatientRepoModel, error)
		GetAll(ctx context.Context, currentUser UserRepoModel) ([]PatientRepoModel, error)
		GetOne(ctx context.Context, currentUser UserRepoModel, patientID uuid.UUID) (PatientRepoModel, error)
		LinkPatient(ctx context.Context, user UserRepoModel, patientID string) error
		GetAllByUser(ctx context.Context, user UserRepoModel) ([]PatientAccountInfo, error)
	}

	PatientRepository interface {
		SelectAIGeneratedReportsByPatient(ctx context.Context, patientID uuid.UUID) ([]ReportForPatient, error)
		GetSelfPatientIDOfUser(ctx context.Context, userID int) (uuid.UUID, error)
		SelectApprovedReportsByUser(ctx context.Context, userID int) ([]ReportForPatient, error)
		InsertPatient(ctx context.Context, patient PatientRepoModel) (PatientRepoModel, error)
		SelectAll(ctx context.Context, userID int) ([]PatientRepoModel, error)
		SelectOne(ctx context.Context, userID int, patientID uuid.UUID) (PatientRepoModel, error)
		SelectAllByUser(ctx context.Context, userID int) ([]PatientAccountInfo, error)
		LinkPatient(ctx context.Context, user UserRepoModel, patientID uuid.UUID) error
	}
)
