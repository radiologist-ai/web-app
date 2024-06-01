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
		GetReportsByPatient(ctx context.Context, patientID uuid.UUID) ([]ReportModel, error)
		UpdateReport(ctx context.Context, id int, opts ...PatchOpt) error
		GetOne(ctx context.Context, id int) (ReportModel, error)
		GenerateReport(ctx context.Context, patientID uuid.UUID, photo io.Reader, ext string) (ReportModel, error)
	}
	ReportRepository interface {
		GetReportsByPatient(ctx context.Context, patientID uuid.UUID) ([]ReportModel, error)
		CreateReport(ctx context.Context, patientID uuid.UUID, imagePath, reportText string, approved bool) (createdModel ReportModel, err error)
		PatchReport(ctx context.Context, id int, opts ...PatchOpt) error
		SelectReport(ctx context.Context, id int) (ReportModel, error)
	}

	ReportPatchForm struct {
		ID         int    `db:"id"`
		ReportText string `db:"report_text"`
		Approved   bool   `db:"approved"`
	}
	ReportModel struct {
		ID         int       `db:"id"`
		PatientID  uuid.UUID `db:"patient_id"`
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

const (
	XrayWidth  = 256
	XrayHeight = 256
)

func WithReportText(s string) PatchOpt {
	return func(conf *PatchConf) {
		conf.ReportText = &s
	}
}

func WithReportApproved(b bool) PatchOpt {
	return func(conf *PatchConf) {
		conf.Approved = &b
	}
}
