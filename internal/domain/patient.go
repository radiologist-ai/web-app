package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type PatientRepoModel struct {
	ID                uuid.UUID `db:"id"`
	UserID            *int      `db:"user_id"`
	CreatorID         int       `db:"creator_id"`
	Name              string    `db:"name"`
	PatientIdentifier string    `db:"patient_identifier"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

type PatientService interface {
	CreatePatient(ctx context.Context, creator UserRepoModel, form PatientRepoModel) (PatientRepoModel, error)
	GetAll(ctx context.Context, currentUser UserRepoModel) ([]PatientRepoModel, error)
}

type PatientRepository interface {
	InsertPatient(ctx context.Context, patient PatientRepoModel) (PatientRepoModel, error)
	SelectAll(ctx context.Context, userID int) ([]PatientRepoModel, error)
}
