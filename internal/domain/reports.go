package domain

import (
	"context"
	"github.com/google/uuid"
	"io"
	"time"
)

type (
	RGen interface {
		GenerateReport(ctx context.Context, link2photo string) (string, error)
		GenerateReportAsync(ctx context.Context, link2photo string, ch chan string, errCh chan error)
	}
	ReportService interface {
		GenerateReport(ctx context.Context, patientID uuid.UUID, photo io.ReadCloser, size int64) (ReportModel, error)
	}
	ReportRepository interface {
		CreateReport(ctx context.Context, patientID int, imagePath, reportText string, approved bool) (createdModel ReportModel, err error)
		PatchReport(ctx context.Context, id int, opts ...PatchOpt) error
	}

	ReportPatchForm struct {
		ID         int    `db:"id"`
		ReportText string `db:"report_text"`
		Approved   bool   `db:"approved"`
	}
	ReportModel struct {
		ID         int       `db:"id"`
		PatientID  int       `db:"patient_id"`
		ImagePath  string    `db:"image_path"`
		ReportText string    `db:"report_text"`
		Approved   bool      `db:"approved"`
		CreatedAt  time.Time `db:"created_at"`
		UpdatedAt  time.Time `db:"updated_at"`
	}
	PatchConf struct {
		ReportText *string
		Approved   *bool
	}

	PatchOpt func(*PatchConf)
)
